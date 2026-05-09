//! Minimal HTTP-over-UDS client for the Firecracker API.

use std::path::{Path, PathBuf};

use anyhow::{Context, Result, anyhow, bail};
use serde::{Deserialize, Serialize, de::DeserializeOwned};
use tokio::{
    io::{AsyncReadExt, AsyncWriteExt},
    net::UnixStream,
};

#[derive(Debug, Clone)]
pub struct ApiClient {
    socket_path: PathBuf,
}

impl ApiClient {
    pub fn new(socket_path: impl Into<PathBuf>) -> Self {
        Self {
            socket_path: socket_path.into(),
        }
    }

    pub async fn put_json<T>(&self, path: &str, body: &T) -> Result<()>
    where
        T: Serialize + ?Sized,
    {
        self.request("PUT", path, Some(serde_json::to_vec(body)?))
            .await?;
        Ok(())
    }

    pub async fn get_json<T>(&self, path: &str) -> Result<T>
    where
        T: DeserializeOwned,
    {
        let response = self.request("GET", path, None).await?;
        serde_json::from_slice(&response)
            .with_context(|| format!("decode Firecracker response from {path}"))
    }

    pub async fn configure_machine(&self, config: &MachineConfiguration) -> Result<()> {
        self.put_json("/machine-config", config).await
    }

    pub async fn set_boot_source(&self, boot_source: &BootSource) -> Result<()> {
        self.put_json("/boot-source", boot_source).await
    }

    pub async fn put_drive(&self, drive: &Drive) -> Result<()> {
        self.put_json(&format!("/drives/{}", drive.drive_id), drive)
            .await
    }

    pub async fn put_network_interface(&self, interface: &NetworkInterface) -> Result<()> {
        self.put_json(
            &format!("/network-interfaces/{}", interface.iface_id),
            interface,
        )
        .await
    }

    pub async fn start_instance(&self) -> Result<()> {
        self.put_json(
            "/actions",
            &Action {
                action_type: ActionType::InstanceStart,
            },
        )
        .await
    }

    pub async fn send_ctrl_alt_del(&self) -> Result<()> {
        self.put_json(
            "/actions",
            &Action {
                action_type: ActionType::SendCtrlAltDel,
            },
        )
        .await
    }

    async fn request(&self, method: &str, path: &str, body: Option<Vec<u8>>) -> Result<Vec<u8>> {
        if !path.starts_with('/') {
            bail!("Firecracker API path must start with '/': {path}");
        }

        let body = body.unwrap_or_default();
        let request = format!(
            "{method} {path} HTTP/1.1\r\n\
             Host: localhost\r\n\
             Accept: application/json\r\n\
             Content-Type: application/json\r\n\
             Content-Length: {}\r\n\
             Connection: close\r\n\
             \r\n",
            body.len()
        );

        let mut stream = UnixStream::connect(&self.socket_path)
            .await
            .with_context(|| {
                format!(
                    "connect Firecracker API socket {}",
                    display_path(&self.socket_path)
                )
            })?;
        stream.write_all(request.as_bytes()).await?;
        stream.write_all(&body).await?;
        stream.shutdown().await?;

        let mut response = Vec::new();
        stream.read_to_end(&mut response).await?;
        parse_http_response(&response).with_context(|| {
            format!(
                "Firecracker API {method} {path} via {}",
                display_path(&self.socket_path)
            )
        })
    }
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub struct MachineConfiguration {
    pub vcpu_count: u8,
    pub mem_size_mib: u32,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub smt: Option<bool>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub track_dirty_pages: Option<bool>,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub struct BootSource {
    pub kernel_image_path: PathBuf,
    pub boot_args: String,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub struct Drive {
    pub drive_id: String,
    pub path_on_host: PathBuf,
    pub is_root_device: bool,
    pub is_read_only: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub struct NetworkInterface {
    pub iface_id: String,
    pub host_dev_name: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub guest_mac: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub struct Action {
    pub action_type: ActionType,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub enum ActionType {
    InstanceStart,
    SendCtrlAltDel,
}

fn parse_http_response(response: &[u8]) -> Result<Vec<u8>> {
    let header_end = response
        .windows(4)
        .position(|window| window == b"\r\n\r\n")
        .ok_or_else(|| anyhow!("malformed HTTP response: missing header terminator"))?;
    let headers = std::str::from_utf8(&response[..header_end])
        .context("malformed HTTP response: headers are not utf-8")?;
    let body = response[header_end + 4..].to_vec();

    let status_line = headers
        .lines()
        .next()
        .ok_or_else(|| anyhow!("malformed HTTP response: missing status line"))?;
    let status = status_line
        .split_whitespace()
        .nth(1)
        .ok_or_else(|| anyhow!("malformed HTTP response: missing status code"))?
        .parse::<u16>()
        .context("malformed HTTP response: invalid status code")?;

    if !(200..300).contains(&status) {
        let body_text = String::from_utf8_lossy(&body);
        bail!("Firecracker API returned HTTP {status}: {body_text}");
    }

    Ok(body)
}

fn display_path(path: &Path) -> String {
    path.display().to_string()
}

#[cfg(test)]
mod tests {
    use super::*;
    use serde_json::json;
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

    static SOCKET_COUNTER: AtomicU64 = AtomicU64::new(0);

    #[tokio::test]
    async fn put_json_sends_firecracker_request_over_unix_socket() {
        let socket_path = unique_socket_path();
        let captured = Arc::new(Mutex::new(String::new()));
        let server = spawn_one_response_server(
            &socket_path,
            "HTTP/1.1 204 No Content\r\nContent-Length: 0\r\n\r\n",
            captured.clone(),
        )
        .await;

        let client = ApiClient::new(&socket_path);
        client
            .configure_machine(&MachineConfiguration {
                vcpu_count: 2,
                mem_size_mib: 1024,
                smt: Some(false),
                track_dirty_pages: None,
            })
            .await
            .expect("PUT should succeed");
        server.await.expect("server task should finish");

        let request = captured.lock().expect("captured request").clone();
        assert!(request.starts_with("PUT /machine-config HTTP/1.1\r\n"));
        assert!(request.contains("Content-Type: application/json\r\n"));
        assert!(request.contains(r#""vcpu_count":2"#));
        assert!(request.contains(r#""mem_size_mib":1024"#));
        assert!(request.contains(r#""smt":false"#));
        assert!(!request.contains("track_dirty_pages"));
    }

    #[tokio::test]
    async fn get_json_decodes_successful_response_body() {
        let socket_path = unique_socket_path();
        let captured = Arc::new(Mutex::new(String::new()));
        let server = spawn_one_response_server(
            &socket_path,
            "HTTP/1.1 200 OK\r\nContent-Length: 36\r\n\r\n{\"vcpu_count\":1,\"mem_size_mib\":512}",
            captured.clone(),
        )
        .await;

        let client = ApiClient::new(&socket_path);
        let config: MachineConfiguration = client
            .get_json("/machine-config")
            .await
            .expect("GET should decode JSON");
        server.await.expect("server task should finish");

        assert_eq!(
            config,
            MachineConfiguration {
                vcpu_count: 1,
                mem_size_mib: 512,
                smt: None,
                track_dirty_pages: None,
            }
        );
        assert!(
            captured
                .lock()
                .expect("captured request")
                .starts_with("GET /machine-config HTTP/1.1\r\n")
        );
    }

    #[tokio::test]
    async fn non_success_status_includes_response_body() {
        let socket_path = unique_socket_path();
        let captured = Arc::new(Mutex::new(String::new()));
        let server = spawn_one_response_server(
            &socket_path,
            "HTTP/1.1 400 Bad Request\r\nContent-Length: 24\r\n\r\n{\"fault\":\"bad request\"}",
            captured,
        )
        .await;

        let client = ApiClient::new(&socket_path);
        let err = client
            .start_instance()
            .await
            .expect_err("HTTP error should fail");
        server.await.expect("server task should finish");

        let err = format!("{err:#}");
        assert!(err.contains("HTTP 400"));
        assert!(err.contains("bad request"));
    }

    #[tokio::test]
    async fn rejects_relative_api_paths() {
        let client = ApiClient::new("/tmp/missing-firecracker.sock");

        let err = client
            .get_json::<serde_json::Value>("machine-config")
            .await
            .expect_err("relative path should fail before connecting");

        assert!(err.to_string().contains("must start with '/'"));
    }

    fn unique_socket_path() -> PathBuf {
        std::env::temp_dir().join(format!(
            "ssh-microvm-api-test-{}-{}.sock",
            std::process::id(),
            SOCKET_COUNTER.fetch_add(1, Ordering::Relaxed)
        ))
    }

    async fn spawn_one_response_server(
        socket_path: &Path,
        response: &'static str,
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
                .write_all(response.as_bytes())
                .await
                .expect("write response");
        })
    }

    #[test]
    fn serializes_common_firecracker_payloads() {
        assert_eq!(
            serde_json::to_value(BootSource {
                kernel_image_path: PathBuf::from("/vm/vmlinux"),
                boot_args: "console=ttyS0".to_string(),
            })
            .expect("serialize boot source"),
            json!({
                "kernel_image_path": "/vm/vmlinux",
                "boot_args": "console=ttyS0",
            })
        );

        assert_eq!(
            serde_json::to_value(Drive {
                drive_id: "rootfs".to_string(),
                path_on_host: PathBuf::from("/vm/rootfs.ext4"),
                is_root_device: true,
                is_read_only: false,
            })
            .expect("serialize drive"),
            json!({
                "drive_id": "rootfs",
                "path_on_host": "/vm/rootfs.ext4",
                "is_root_device": true,
                "is_read_only": false,
            })
        );

        assert_eq!(
            serde_json::to_value(NetworkInterface {
                iface_id: "eth0".to_string(),
                host_dev_name: "tap0".to_string(),
                guest_mac: None,
            })
            .expect("serialize network interface"),
            json!({
                "iface_id": "eth0",
                "host_dev_name": "tap0",
            })
        );

        assert_eq!(
            serde_json::to_value(Action {
                action_type: ActionType::SendCtrlAltDel,
            })
            .expect("serialize action"),
            json!({
                "action_type": "SendCtrlAltDel",
            })
        );
    }
}
