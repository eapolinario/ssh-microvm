//! Inner SSH client used to execute commands inside the guest.

use std::{net::SocketAddr, sync::Arc};

use anyhow::{Context, Result, anyhow};
use russh::{
    ChannelMsg, Disconnect, client,
    keys::{PrivateKeyWithHashAlg, PublicKey, load_secret_key},
};

use crate::{config::Config, lifecycle::BoxFuture};

pub(crate) struct ExecOutput {
    pub stdout: Vec<u8>,
    pub stderr: Vec<u8>,
    pub exit_status: u32,
}

pub(crate) trait ExecProxy: Send + Sync + 'static {
    fn exec(&self, command: Vec<u8>) -> BoxFuture<Result<ExecOutput>>;
}

#[derive(Clone)]
pub(crate) struct GuestSshProxy {
    cfg: Arc<Config>,
    addr: SocketAddr,
}

impl GuestSshProxy {
    pub(crate) fn new(cfg: &Config) -> Self {
        Self::new_with_addr(cfg, SocketAddr::new(cfg.guest_ip, 22))
    }

    pub(crate) fn new_with_addr(cfg: &Config, addr: SocketAddr) -> Self {
        Self {
            cfg: Arc::new(cfg.clone()),
            addr,
        }
    }
}

impl ExecProxy for GuestSshProxy {
    fn exec(&self, command: Vec<u8>) -> BoxFuture<Result<ExecOutput>> {
        let cfg = self.cfg.clone();
        let addr = self.addr;

        Box::pin(async move {
            let key = load_secret_key(&cfg.guest_key, None)
                .with_context(|| format!("load guest SSH key {}", cfg.guest_key.display()))?;
            let mut session =
                client::connect(Arc::new(client::Config::default()), addr, GuestClient)
                    .await
                    .with_context(|| format!("connect to guest sshd at {addr}"))?;
            let authenticated = session
                .authenticate_publickey(
                    cfg.guest_user.clone(),
                    PrivateKeyWithHashAlg::new(
                        Arc::new(key),
                        session
                            .best_supported_rsa_hash()
                            .await
                            .context("query guest SSH server RSA signature support")?
                            .flatten(),
                    ),
                )
                .await
                .context("authenticate to guest sshd with --guest-key")?
                .success();
            if !authenticated {
                return Err(anyhow!("guest SSH public-key authentication failed"));
            }

            let mut channel = session
                .channel_open_session()
                .await
                .context("open guest SSH session channel")?;
            channel
                .exec(true, command)
                .await
                .context("send guest SSH exec request")?;

            let mut stdout = Vec::new();
            let mut stderr = Vec::new();
            let mut exit_status = None;
            while let Some(msg) = channel.wait().await {
                match msg {
                    ChannelMsg::Data { data } => stdout.extend_from_slice(&data),
                    ChannelMsg::ExtendedData { data, .. } => stderr.extend_from_slice(&data),
                    ChannelMsg::ExitStatus { exit_status: code } => exit_status = Some(code),
                    ChannelMsg::Close => break,
                    _ => {}
                }
            }

            session
                .disconnect(Disconnect::ByApplication, "", "en")
                .await
                .context("disconnect guest SSH session")?;

            Ok(ExecOutput {
                stdout,
                stderr,
                exit_status: exit_status
                    .ok_or_else(|| anyhow!("guest SSH exec ended without exit status"))?,
            })
        })
    }
}

struct GuestClient;

impl client::Handler for GuestClient {
    type Error = anyhow::Error;

    async fn check_server_key(&mut self, _server_public_key: &PublicKey) -> Result<bool> {
        Ok(true)
    }
}

#[cfg(test)]
mod tests {
    use std::{
        fs,
        net::{IpAddr, Ipv4Addr},
        path::PathBuf,
        time::Duration,
    };

    use russh::{
        Channel, ChannelId, MethodKind, MethodSet,
        keys::{Algorithm, PrivateKey, key::safe_rng, ssh_key::LineEnding},
        server::{self, Auth, Msg, Server as _, Session},
    };
    use tokio::net::TcpListener;

    use crate::config::Command;

    use super::*;

    #[tokio::test]
    async fn guest_proxy_exec_collects_output_and_status() {
        let client_key =
            PrivateKey::random(&mut safe_rng(), Algorithm::Ed25519).expect("client key");
        let host_key = PrivateKey::random(&mut safe_rng(), Algorithm::Ed25519).expect("host key");
        let key_path = write_test_key(&client_key);
        let listener = TcpListener::bind("127.0.0.1:0").await.expect("bind");
        let addr = listener.local_addr().expect("local addr");
        let server_config = Arc::new(server::Config {
            keys: vec![host_key],
            methods: MethodSet::from(&[MethodKind::PublicKey][..]),
            auth_rejection_time_initial: Some(Duration::from_millis(0)),
            ..Default::default()
        });
        let mut server = InnerServer {
            accepted: client_key.public_key().clone(),
        };

        let server_task = tokio::spawn(async move {
            let (socket, peer_addr) = listener.accept().await.expect("accept");
            let handler = server.new_client(Some(peer_addr));
            server::run_stream(server_config, socket, handler)
                .await
                .expect("run inner ssh server");
        });

        let mut cfg = test_config();
        cfg.guest_key = key_path.clone();
        cfg.guest_user = "guest".to_string();
        let proxy = GuestSshProxy::new_with_addr(&cfg, addr);

        let output = proxy
            .exec(b"echo hello".to_vec())
            .await
            .expect("proxy exec");

        assert_eq!(output.stdout, b"ran: echo hello\n");
        assert_eq!(output.stderr, b"warn\n");
        assert_eq!(output.exit_status, 7);
        server_task.await.expect("server task");
        fs::remove_dir_all(key_path.parent().unwrap()).expect("remove key dir");
    }

    struct InnerServer {
        accepted: PublicKey,
    }

    impl server::Server for InnerServer {
        type Handler = InnerSession;

        fn new_client(&mut self, _peer_addr: Option<std::net::SocketAddr>) -> Self::Handler {
            InnerSession {
                accepted: self.accepted.clone(),
            }
        }
    }

    struct InnerSession {
        accepted: PublicKey,
    }

    impl server::Handler for InnerSession {
        type Error = anyhow::Error;

        async fn auth_publickey_offered(
            &mut self,
            _user: &str,
            public_key: &PublicKey,
        ) -> Result<Auth, Self::Error> {
            Ok(if public_key == &self.accepted {
                Auth::Accept
            } else {
                Auth::reject()
            })
        }

        async fn auth_publickey(
            &mut self,
            _user: &str,
            public_key: &PublicKey,
        ) -> Result<Auth, Self::Error> {
            Ok(if public_key == &self.accepted {
                Auth::Accept
            } else {
                Auth::reject()
            })
        }

        async fn channel_open_session(
            &mut self,
            _channel: Channel<Msg>,
            _session: &mut Session,
        ) -> Result<bool, Self::Error> {
            Ok(true)
        }

        async fn exec_request(
            &mut self,
            channel: ChannelId,
            data: &[u8],
            session: &mut Session,
        ) -> Result<(), Self::Error> {
            session.channel_success(channel)?;
            session.data(channel, format!("ran: {}\n", String::from_utf8_lossy(data)))?;
            session.extended_data(channel, 1, "warn\n")?;
            session.exit_status_request(channel, 7)?;
            session.eof(channel)?;
            session.close(channel)?;
            Ok(())
        }
    }

    fn write_test_key(key: &PrivateKey) -> PathBuf {
        let dir = std::env::temp_dir().join(format!(
            "ssh-microvm-guest-key-{}-{}",
            std::process::id(),
            std::time::SystemTime::now()
                .duration_since(std::time::UNIX_EPOCH)
                .expect("system time after epoch")
                .as_nanos()
        ));
        fs::create_dir_all(&dir).expect("create key dir");
        let path = dir.join("guest_key");
        fs::write(
            &path,
            key.to_openssh(LineEnding::LF)
                .expect("encode private key")
                .as_bytes(),
        )
        .expect("write key");
        path
    }

    fn test_config() -> Config {
        Config {
            dry_boot: false,
            command: None::<Command>,
            listen: "127.0.0.1:0".parse().unwrap(),
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
            guest_ip: IpAddr::V4(Ipv4Addr::LOCALHOST),
            host_ip: IpAddr::V4(Ipv4Addr::LOCALHOST),
            tap_name: "tap0".to_string(),
            boot_timeout: Duration::from_secs(1),
            grace_stop: Duration::from_secs(1),
        }
    }
}
