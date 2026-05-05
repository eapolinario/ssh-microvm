//! Firecracker process lifecycle management.

use std::{
    path::{Path, PathBuf},
    process::Stdio,
    time::Duration,
};

use anyhow::{Context, Result, bail};
use tokio::{
    net::UnixStream,
    process::{Child, Command},
    time::{Instant, sleep, timeout},
};

use crate::{
    api::ApiClient,
    boot::{api_socket_path, configure_and_start, wait_for_guest_sshd},
    config::Config,
};

const FIRECRACKER_API_POLL_INTERVAL: Duration = Duration::from_millis(50);
const FIRECRACKER_API_CONNECT_ATTEMPT_TIMEOUT: Duration = Duration::from_millis(200);

pub struct VmHandle {
    child: Child,
    client: ApiClient,
    state_dir: PathBuf,
    socket_path: PathBuf,
    grace_stop: Duration,
}

impl VmHandle {
    pub async fn boot(cfg: &Config) -> Result<Self> {
        let mut vm = Self::spawn(cfg).await?;

        if let Err(err) = vm.wait_until_ready(cfg).await {
            let cleanup_result = vm.kill_and_cleanup().await;
            if let Err(cleanup_err) = cleanup_result {
                return Err(err).with_context(|| {
                    format!("also failed to clean up unsuccessful VM boot: {cleanup_err:#}")
                });
            }
            return Err(err);
        }

        Ok(vm)
    }

    pub(crate) async fn spawn(cfg: &Config) -> Result<Self> {
        tokio::fs::create_dir_all(&cfg.state_dir)
            .await
            .with_context(|| format!("create state directory {}", cfg.state_dir.display()))?;

        let socket_path = api_socket_path(cfg);
        remove_stale_socket(&socket_path).await?;

        let child = spawn_firecracker(cfg, &socket_path)?;
        tracing::info!(
            firecracker = %cfg.firecracker.display(),
            api_socket = %socket_path.display(),
            "spawned Firecracker"
        );

        Ok(Self {
            child,
            client: ApiClient::new(&socket_path),
            state_dir: cfg.state_dir.clone(),
            socket_path,
            grace_stop: cfg.grace_stop,
        })
    }

    pub(crate) async fn wait_until_ready(&mut self, cfg: &Config) -> Result<()> {
        wait_for_api_socket(&self.socket_path, cfg.boot_timeout).await?;
        configure_and_start(&self.client, cfg).await?;
        wait_for_guest_sshd(cfg.guest_ip, cfg.boot_timeout).await?;
        tracing::info!(guest_ip = %cfg.guest_ip, "guest sshd is reachable");
        Ok(())
    }

    pub async fn shutdown(mut self) -> Result<()> {
        let shutdown_result = self
            .client
            .send_ctrl_alt_del()
            .await
            .context("request Firecracker guest shutdown with SendCtrlAltDel");

        let stop_result = if shutdown_result.is_ok() {
            self.wait_for_exit_or_kill().await
        } else {
            self.kill_child().await
        };
        let cleanup_result = cleanup_state_dir(&self.state_dir).await;

        shutdown_result?;
        stop_result?;
        cleanup_result?;
        Ok(())
    }

    pub(crate) async fn kill_and_cleanup(mut self) -> Result<()> {
        let stop_result = self.kill_child().await;
        let cleanup_result = cleanup_state_dir(&self.state_dir).await;

        stop_result?;
        cleanup_result?;
        Ok(())
    }

    async fn wait_for_exit_or_kill(&mut self) -> Result<()> {
        if self
            .child
            .try_wait()
            .context("check Firecracker process status")?
            .is_some()
        {
            return Ok(());
        }

        match timeout(self.grace_stop, self.child.wait()).await {
            Ok(result) => {
                result.context("wait for Firecracker process to exit after SendCtrlAltDel")?;
                Ok(())
            }
            Err(_elapsed) => self.kill_child().await,
        }
    }

    async fn kill_child(&mut self) -> Result<()> {
        if self
            .child
            .try_wait()
            .context("check Firecracker process status")?
            .is_some()
        {
            return Ok(());
        }

        self.child
            .kill()
            .await
            .context("kill Firecracker process")?;
        Ok(())
    }
}

fn spawn_firecracker(cfg: &Config, socket_path: &Path) -> Result<Child> {
    Command::new(&cfg.firecracker)
        .arg("--api-sock")
        .arg(socket_path)
        .stdin(Stdio::null())
        .stdout(Stdio::inherit())
        .stderr(Stdio::inherit())
        .spawn()
        .with_context(|| format!("spawn Firecracker {}", cfg.firecracker.display()))
}

async fn wait_for_api_socket(socket_path: &Path, boot_timeout: Duration) -> Result<()> {
    let deadline = Instant::now() + boot_timeout;

    loop {
        let last_error = match timeout(
            FIRECRACKER_API_CONNECT_ATTEMPT_TIMEOUT,
            UnixStream::connect(socket_path),
        )
        .await
        {
            Ok(Ok(_stream)) => return Ok(()),
            Ok(Err(err)) => err.to_string(),
            Err(_elapsed) => "connection attempt timed out".to_string(),
        };

        let now = Instant::now();
        if now >= deadline {
            bail!(
                "Firecracker API socket {} did not become reachable within {:?}: {last_error}",
                socket_path.display(),
                boot_timeout
            );
        }

        sleep(std::cmp::min(FIRECRACKER_API_POLL_INTERVAL, deadline - now)).await;
    }
}

async fn remove_stale_socket(socket_path: &Path) -> Result<()> {
    match tokio::fs::remove_file(socket_path).await {
        Ok(()) => Ok(()),
        Err(err) if err.kind() == std::io::ErrorKind::NotFound => Ok(()),
        Err(err) => Err(err)
            .with_context(|| format!("remove Firecracker API socket {}", socket_path.display())),
    }
}

async fn cleanup_state_dir(state_dir: &Path) -> Result<()> {
    match tokio::fs::remove_dir_all(state_dir).await {
        Ok(()) => Ok(()),
        Err(err) if err.kind() == std::io::ErrorKind::NotFound => Ok(()),
        Err(err) => Err(err)
            .with_context(|| format!("remove Firecracker state directory {}", state_dir.display())),
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use clap::Parser;
    use std::{
        fs,
        sync::{
            Arc, Mutex,
            atomic::{AtomicU64, Ordering},
        },
    };
    use tokio::{
        io::{AsyncReadExt, AsyncWriteExt},
        net::UnixListener,
    };

    static STATE_COUNTER: AtomicU64 = AtomicU64::new(0);

    #[tokio::test]
    async fn shutdown_sends_ctrl_alt_del_kills_after_grace_and_removes_state_dir() {
        let state_dir = unique_state_dir();
        fs::create_dir_all(&state_dir).expect("create state dir");
        let socket_path = state_dir.join("firecracker.sock");
        let request = Arc::new(Mutex::new(String::new()));
        let server = spawn_one_response_server(&socket_path, request.clone()).await;
        let child = Command::new("sleep")
            .arg("60")
            .spawn()
            .expect("spawn sleep process");

        let vm = VmHandle {
            child,
            client: ApiClient::new(&socket_path),
            state_dir: state_dir.clone(),
            socket_path,
            grace_stop: Duration::from_millis(10),
        };

        vm.shutdown().await.expect("shutdown should succeed");
        server.await.expect("server task should finish");

        assert!(
            request
                .lock()
                .expect("captured request")
                .starts_with("PUT /actions HTTP/1.1\r\n")
        );
        assert!(
            request
                .lock()
                .expect("captured request")
                .contains(r#""action_type":"SendCtrlAltDel""#)
        );
        assert!(
            !state_dir.exists(),
            "shutdown should remove the VM state directory"
        );
    }

    #[tokio::test]
    async fn shutdown_force_kills_and_cleans_up_when_ctrl_alt_del_fails() {
        let state_dir = unique_state_dir();
        fs::create_dir_all(&state_dir).expect("create state dir");
        let socket_path = state_dir.join("missing.sock");
        let child = Command::new("sleep")
            .arg("60")
            .spawn()
            .expect("spawn sleep process");

        let vm = VmHandle {
            child,
            client: ApiClient::new(&socket_path),
            state_dir: state_dir.clone(),
            socket_path,
            grace_stop: Duration::from_secs(1),
        };

        let err = vm
            .shutdown()
            .await
            .expect_err("missing API socket should surface as shutdown error");

        assert!(format!("{err:#}").contains("SendCtrlAltDel"));
        assert!(
            !state_dir.exists(),
            "failed graceful shutdown should still remove state directory"
        );
    }

    #[tokio::test]
    async fn wait_for_api_socket_returns_when_socket_accepts_connections() {
        let state_dir = unique_state_dir();
        fs::create_dir_all(&state_dir).expect("create state dir");
        let socket_path = state_dir.join("firecracker.sock");
        let _ = fs::remove_file(&socket_path);
        let listener = UnixListener::bind(&socket_path).expect("bind unix listener");
        let server = tokio::spawn(async move {
            let _ = listener.accept().await.expect("accept connection");
        });

        wait_for_api_socket(&socket_path, Duration::from_secs(1))
            .await
            .expect("reachable API socket should succeed");
        server.await.expect("server task should finish");
        fs::remove_dir_all(&state_dir).expect("remove state dir");
    }

    #[test]
    fn api_socket_path_lives_under_state_dir() {
        let cfg = test_config();

        assert_eq!(
            api_socket_path(&cfg),
            PathBuf::from("/tmp/ssh-microvm/firecracker.sock")
        );
    }

    fn test_config() -> Config {
        Config::try_parse_from([
            "ssh-microvm",
            "--kernel",
            "artifacts/vmlinux.bin",
            "--rootfs",
            "artifacts/ubuntu.ext4",
            "--guest-key",
            "artifacts/ubuntu.id_rsa",
            "--accept-any-key",
            "--vcpu",
            "2",
            "--mem",
            "1024",
            "--boot-args",
            "console=ttyS0",
            "--tap-name",
            "tap-test",
        ])
        .expect("valid test config")
    }

    fn unique_state_dir() -> PathBuf {
        std::env::temp_dir().join(format!(
            "ssh-microvm-firecracker-test-{}-{}",
            std::process::id(),
            STATE_COUNTER.fetch_add(1, Ordering::Relaxed)
        ))
    }

    async fn spawn_one_response_server(
        socket_path: &Path,
        captured: Arc<Mutex<String>>,
    ) -> tokio::task::JoinHandle<()> {
        let _ = fs::remove_file(socket_path);
        let listener = UnixListener::bind(socket_path).expect("bind unix listener");

        tokio::spawn(async move {
            let (mut stream, _) = listener.accept().await.expect("accept client");
            let mut request = Vec::new();
            stream
                .read_to_end(&mut request)
                .await
                .expect("read request");
            *captured.lock().expect("captured request") =
                String::from_utf8(request).expect("request is utf-8");
            stream
                .write_all(b"HTTP/1.1 204 No Content\r\nContent-Length: 0\r\n\r\n")
                .await
                .expect("write response");
        })
    }
}
