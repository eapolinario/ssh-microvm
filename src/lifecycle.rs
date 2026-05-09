//! Single-VM lifecycle gate for SSH sessions.

use std::{future::Future, pin::Pin, sync::Arc};

use anyhow::Result;
use tokio::sync::{Mutex, OwnedMutexGuard};

use crate::{config::Config, firecracker::VmHandle};

pub(crate) type BoxFuture<T> = Pin<Box<dyn Future<Output = T> + Send>>;

pub(crate) trait RunningVm: Send + 'static {
    fn shutdown(self: Box<Self>) -> BoxFuture<Result<()>>;
}

pub(crate) trait VmSpawner: Send + Sync + 'static {
    fn boot(&self, cfg: Config) -> BoxFuture<Result<Box<dyn RunningVm>>>;
}

pub struct VmLifecycle {
    cfg: Arc<Config>,
    gate: Arc<Mutex<()>>,
    spawner: Arc<dyn VmSpawner>,
}

impl VmLifecycle {
    pub fn new(cfg: &Config) -> Self {
        Self::new_with_spawner(cfg, Arc::new(FirecrackerSpawner))
    }

    pub(crate) fn new_with_spawner(cfg: &Config, spawner: Arc<dyn VmSpawner>) -> Self {
        Self {
            cfg: Arc::new(cfg.clone()),
            gate: Arc::new(Mutex::new(())),
            spawner,
        }
    }

    pub async fn acquire(&self) -> Result<VmLease> {
        let guard = self.gate.clone().lock_owned().await;
        let vm = self.spawner.boot((*self.cfg).clone()).await?;
        Ok(VmLease {
            vm: Some(vm),
            guard: Some(guard),
        })
    }
}

impl Clone for VmLifecycle {
    fn clone(&self) -> Self {
        Self {
            cfg: self.cfg.clone(),
            gate: self.gate.clone(),
            spawner: self.spawner.clone(),
        }
    }
}

pub struct VmLease {
    vm: Option<Box<dyn RunningVm>>,
    guard: Option<OwnedMutexGuard<()>>,
}

impl Drop for VmLease {
    fn drop(&mut self) {
        let Some(vm) = self.vm.take() else {
            return;
        };
        let guard = self.guard.take();

        tokio::spawn(async move {
            if let Err(err) = vm.shutdown().await {
                tracing::warn!(error = ?err, "failed to shut down VM after SSH disconnect");
            }
            drop(guard);
        });
    }
}

struct FirecrackerSpawner;

impl VmSpawner for FirecrackerSpawner {
    fn boot(&self, cfg: Config) -> BoxFuture<Result<Box<dyn RunningVm>>> {
        Box::pin(async move {
            let vm = VmHandle::boot(&cfg).await?;
            Ok(Box::new(FirecrackerVm(vm)) as Box<dyn RunningVm>)
        })
    }
}

struct FirecrackerVm(VmHandle);

impl RunningVm for FirecrackerVm {
    fn shutdown(self: Box<Self>) -> BoxFuture<Result<()>> {
        Box::pin(async move {
            let FirecrackerVm(vm) = *self;
            vm.shutdown().await
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::config::Config;
    use std::{
        net::{IpAddr, Ipv4Addr, SocketAddr},
        path::PathBuf,
        sync::{
            Arc,
            atomic::{AtomicUsize, Ordering},
        },
        time::Duration,
    };
    use tokio::{sync::Notify, time::timeout};

    #[tokio::test]
    async fn lease_holds_global_gate_until_shutdown_finishes() {
        let spawner = Arc::new(FakeSpawner::default());
        let lifecycle = VmLifecycle::new_with_spawner(&test_config(), spawner.clone());
        let first = lifecycle.acquire().await.expect("first VM boots");

        drop(first);
        spawner.first_shutdown_started.notified().await;

        let second_acquire = tokio::spawn({
            let lifecycle = lifecycle.clone();
            async move { lifecycle.acquire().await }
        });
        timeout(
            Duration::from_millis(20),
            spawner.second_boot_started.notified(),
        )
        .await
        .expect_err("second boot should wait for first shutdown to release the gate");

        spawner.first_shutdown_release.notify_one();
        let second = second_acquire
            .await
            .expect("second acquire task should finish")
            .expect("second VM boots after first shutdown");
        assert_eq!(spawner.boots.load(Ordering::SeqCst), 2);
        drop(second);
    }

    #[derive(Default)]
    struct FakeSpawner {
        boots: Arc<AtomicUsize>,
        first_shutdown_started: Arc<Notify>,
        first_shutdown_release: Arc<Notify>,
        second_boot_started: Arc<Notify>,
    }

    impl VmSpawner for FakeSpawner {
        fn boot(&self, _cfg: Config) -> BoxFuture<Result<Box<dyn RunningVm>>> {
            let boot_number = self.boots.fetch_add(1, Ordering::SeqCst) + 1;
            if boot_number == 2 {
                self.second_boot_started.notify_one();
            }
            let first_shutdown_started = self.first_shutdown_started.clone();
            let first_shutdown_release = self.first_shutdown_release.clone();
            Box::pin(async move {
                Ok(Box::new(FakeVm {
                    boot_number,
                    first_shutdown_started,
                    first_shutdown_release,
                }) as Box<dyn RunningVm>)
            })
        }
    }

    struct FakeVm {
        boot_number: usize,
        first_shutdown_started: Arc<Notify>,
        first_shutdown_release: Arc<Notify>,
    }

    impl RunningVm for FakeVm {
        fn shutdown(self: Box<Self>) -> BoxFuture<Result<()>> {
            Box::pin(async move {
                if self.boot_number == 1 {
                    self.first_shutdown_started.notify_one();
                    self.first_shutdown_release.notified().await;
                }
                Ok(())
            })
        }
    }

    fn test_config() -> Config {
        Config {
            dry_boot: false,
            command: None,
            listen: SocketAddr::new(IpAddr::V4(Ipv4Addr::LOCALHOST), 0),
            kernel: PathBuf::from("kernel"),
            rootfs: PathBuf::from("rootfs"),
            state_dir: PathBuf::from("state"),
            host_key: None,
            authorized_keys: None,
            accept_any_key: true,
            firecracker: PathBuf::from("firecracker"),
            vcpu: 1,
            mem_mib: 512,
            boot_args: String::new(),
            guest_user: "root".to_string(),
            guest_key: PathBuf::from("guest_key"),
            guest_ip: IpAddr::V4(Ipv4Addr::new(172, 16, 0, 2)),
            host_ip: IpAddr::V4(Ipv4Addr::new(172, 16, 0, 1)),
            tap_name: "tap0".to_string(),
            boot_timeout: Duration::from_secs(1),
            grace_stop: Duration::from_secs(1),
        }
    }
}
