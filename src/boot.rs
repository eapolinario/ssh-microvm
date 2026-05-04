//! Firecracker boot helpers shared by CLI boot flows and later VM lifecycle code.

use std::{
    net::{IpAddr, SocketAddr},
    time::Duration,
};

use anyhow::{Result, bail};
use tokio::{
    net::TcpStream,
    time::{Instant, sleep, timeout},
};

use crate::{
    api::{ApiClient, BootSource, Drive, MachineConfiguration, NetworkInterface},
    config::Config,
};

const GUEST_SSH_PORT: u16 = 22;
const SSH_CONNECT_ATTEMPT_TIMEOUT: Duration = Duration::from_millis(200);
const SSH_POLL_INTERVAL: Duration = Duration::from_millis(100);

pub async fn configure_and_start(client: &ApiClient, cfg: &Config) -> Result<()> {
    client
        .configure_machine(&MachineConfiguration {
            vcpu_count: cfg.vcpu,
            mem_size_mib: cfg.mem_mib,
            smt: Some(false),
            track_dirty_pages: None,
        })
        .await?;
    client
        .set_boot_source(&BootSource {
            kernel_image_path: cfg.kernel.clone(),
            boot_args: cfg.boot_args.clone(),
        })
        .await?;
    client
        .put_drive(&Drive {
            drive_id: "rootfs".to_string(),
            path_on_host: cfg.rootfs.clone(),
            is_root_device: true,
            is_read_only: false,
        })
        .await?;
    client
        .put_network_interface(&NetworkInterface {
            iface_id: "eth0".to_string(),
            host_dev_name: cfg.tap_name.clone(),
            guest_mac: None,
        })
        .await?;
    client.start_instance().await
}

pub async fn wait_for_guest_sshd(guest_ip: IpAddr, boot_timeout: Duration) -> Result<()> {
    wait_for_tcp(SocketAddr::new(guest_ip, GUEST_SSH_PORT), boot_timeout).await
}

async fn wait_for_tcp(addr: SocketAddr, boot_timeout: Duration) -> Result<()> {
    let deadline = Instant::now() + boot_timeout;

    loop {
        let last_error = match timeout(SSH_CONNECT_ATTEMPT_TIMEOUT, TcpStream::connect(addr)).await
        {
            Ok(Ok(_stream)) => return Ok(()),
            Ok(Err(err)) => err.to_string(),
            Err(_elapsed) => "connection attempt timed out".to_string(),
        };

        let now = Instant::now();
        if now >= deadline {
            bail!(
                "guest sshd at {addr} did not become reachable within {:?}: {last_error}",
                boot_timeout
            );
        }

        sleep(std::cmp::min(SSH_POLL_INTERVAL, deadline - now)).await;
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use clap::Parser;
    use std::{
        fs,
        path::{Path, PathBuf},
        sync::{
            Arc, Mutex,
            atomic::{AtomicU64, Ordering},
        },
    };
    use tokio::{
        io::{AsyncReadExt, AsyncWriteExt},
        net::{TcpListener, UnixListener},
    };

    static SOCKET_COUNTER: AtomicU64 = AtomicU64::new(0);

    #[tokio::test]
    async fn configure_and_start_sends_firecracker_boot_sequence() {
        let socket_path = unique_socket_path();
        let requests = Arc::new(Mutex::new(Vec::new()));
        let server = spawn_response_server(&socket_path, 5, requests.clone()).await;

        let client = ApiClient::new(&socket_path);
        configure_and_start(&client, &test_config())
            .await
            .expect("configure and start should succeed");
        server.await.expect("server task should finish");

        let requests = requests.lock().expect("captured requests");
        let request_lines = requests
            .iter()
            .map(|request| request.lines().next().expect("request line"))
            .collect::<Vec<_>>();
        assert_eq!(
            request_lines,
            [
                "PUT /machine-config HTTP/1.1",
                "PUT /boot-source HTTP/1.1",
                "PUT /drives/rootfs HTTP/1.1",
                "PUT /network-interfaces/eth0 HTTP/1.1",
                "PUT /actions HTTP/1.1",
            ]
        );

        let bodies = requests.join("\n");
        assert!(bodies.contains(r#""vcpu_count":2"#));
        assert!(bodies.contains(r#""mem_size_mib":1024"#));
        assert!(bodies.contains(r#""kernel_image_path":"artifacts/vmlinux.bin""#));
        assert!(bodies.contains(r#""boot_args":"console=ttyS0""#));
        assert!(bodies.contains(r#""path_on_host":"artifacts/ubuntu.ext4""#));
        assert!(bodies.contains(r#""host_dev_name":"tap-test""#));
        assert!(bodies.contains(r#""action_type":"InstanceStart""#));
    }

    #[tokio::test]
    async fn wait_for_tcp_returns_when_port_accepts_connections() {
        let listener = TcpListener::bind("127.0.0.1:0")
            .await
            .expect("bind test listener");
        let addr = listener.local_addr().expect("listener address");
        let server = tokio::spawn(async move {
            let _ = listener.accept().await.expect("accept connection");
        });

        wait_for_tcp(addr, Duration::from_secs(1))
            .await
            .expect("reachable port should succeed");
        server.await.expect("server task should finish");
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

    fn unique_socket_path() -> PathBuf {
        std::env::temp_dir().join(format!(
            "ssh-microvm-boot-test-{}-{}.sock",
            std::process::id(),
            SOCKET_COUNTER.fetch_add(1, Ordering::Relaxed)
        ))
    }

    async fn spawn_response_server(
        socket_path: &Path,
        expected_requests: usize,
        requests: Arc<Mutex<Vec<String>>>,
    ) -> tokio::task::JoinHandle<()> {
        let _ = fs::remove_file(socket_path);
        let listener = UnixListener::bind(socket_path).expect("bind unix listener");

        tokio::spawn(async move {
            for _ in 0..expected_requests {
                let (mut stream, _) = listener.accept().await.expect("accept client");
                let mut request = Vec::new();
                stream
                    .read_to_end(&mut request)
                    .await
                    .expect("read request");
                requests
                    .lock()
                    .expect("captured requests")
                    .push(String::from_utf8(request).expect("request is utf-8"));
                stream
                    .write_all(b"HTTP/1.1 204 No Content\r\nContent-Length: 0\r\n\r\n")
                    .await
                    .expect("write response");
            }
        })
    }
}
