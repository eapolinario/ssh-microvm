//! Inner SSH client used to execute commands inside the guest.

use std::{io::Cursor, net::SocketAddr, sync::Arc};

use anyhow::{Context, Result, anyhow};
use russh::{
    ChannelMsg, Disconnect, Pty, client,
    keys::{PrivateKeyWithHashAlg, PublicKey, load_secret_key},
};
use tokio::sync::mpsc;

use crate::{config::Config, lifecycle::BoxFuture};

pub(crate) struct ExecOutput {
    pub stdout: Vec<u8>,
    pub stderr: Vec<u8>,
    pub exit_status: u32,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub(crate) struct TerminalSize {
    pub col_width: u32,
    pub row_height: u32,
    pub pix_width: u32,
    pub pix_height: u32,
}

impl TerminalSize {
    pub(crate) fn new(col_width: u32, row_height: u32, pix_width: u32, pix_height: u32) -> Self {
        Self {
            col_width,
            row_height,
            pix_width,
            pix_height,
        }
    }
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub(crate) struct PtyRequest {
    pub term: String,
    pub size: TerminalSize,
    pub modes: Vec<(Pty, u32)>,
}

impl PtyRequest {
    pub(crate) fn new(term: String, size: TerminalSize, modes: Vec<(Pty, u32)>) -> Self {
        Self { term, size, modes }
    }

    pub(crate) fn resize(&mut self, size: TerminalSize) {
        self.size = size;
    }
}

pub(crate) trait ExecProxy: Send + Sync + 'static {
    fn exec(&self, command: Vec<u8>, pty: Option<PtyRequest>) -> BoxFuture<Result<ExecOutput>>;
}

pub(crate) enum ShellEvent {
    Stdout(Vec<u8>),
    Stderr(Vec<u8>),
    ExitStatus(u32),
    Eof,
    Close,
}

pub(crate) struct OpenedShell {
    pub session: Box<dyn ShellSession>,
    pub events: mpsc::Receiver<ShellEvent>,
}

pub(crate) trait ShellSession: Send {
    fn data(&mut self, data: Vec<u8>) -> BoxFuture<Result<()>>;
    fn resize(&mut self, size: TerminalSize) -> BoxFuture<Result<()>>;
    fn eof(&mut self) -> BoxFuture<Result<()>>;
    fn close(&mut self) -> BoxFuture<Result<()>>;
}

pub(crate) trait ShellProxy: Send + Sync + 'static {
    fn open_shell(&self, pty: Option<PtyRequest>) -> BoxFuture<Result<OpenedShell>>;
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
    fn exec(&self, command: Vec<u8>, pty: Option<PtyRequest>) -> BoxFuture<Result<ExecOutput>> {
        let cfg = self.cfg.clone();
        let addr = self.addr;

        Box::pin(async move {
            let session = connect_guest(&cfg, addr).await?;

            let mut channel = session
                .channel_open_session()
                .await
                .context("open guest SSH session channel")?;
            if let Some(pty) = pty {
                channel
                    .request_pty(
                        true,
                        &pty.term,
                        pty.size.col_width,
                        pty.size.row_height,
                        pty.size.pix_width,
                        pty.size.pix_height,
                        &pty.modes,
                    )
                    .await
                    .context("send guest SSH pty request")?;
            }
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

impl ShellProxy for GuestSshProxy {
    fn open_shell(&self, pty: Option<PtyRequest>) -> BoxFuture<Result<OpenedShell>> {
        let cfg = self.cfg.clone();
        let addr = self.addr;

        Box::pin(async move {
            let session = connect_guest(&cfg, addr).await?;
            let channel = session
                .channel_open_session()
                .await
                .context("open guest SSH shell channel")?;
            if let Some(pty) = pty {
                channel
                    .request_pty(
                        true,
                        &pty.term,
                        pty.size.col_width,
                        pty.size.row_height,
                        pty.size.pix_width,
                        pty.size.pix_height,
                        &pty.modes,
                    )
                    .await
                    .context("send guest SSH shell pty request")?;
            }
            channel
                .request_shell(true)
                .await
                .context("send guest SSH shell request")?;

            let (mut read_half, write_half) = channel.split();
            let (events_tx, events_rx) = mpsc::channel(32);
            let (commands_tx, mut commands_rx) = mpsc::channel(32);

            tokio::spawn(async move {
                while let Some(msg) = read_half.wait().await {
                    let event = match msg {
                        ChannelMsg::Data { data } => ShellEvent::Stdout(data.to_vec()),
                        ChannelMsg::ExtendedData { data, .. } => ShellEvent::Stderr(data.to_vec()),
                        ChannelMsg::ExitStatus { exit_status } => {
                            ShellEvent::ExitStatus(exit_status)
                        }
                        ChannelMsg::Eof => ShellEvent::Eof,
                        ChannelMsg::Close => ShellEvent::Close,
                        _ => continue,
                    };
                    let close = matches!(event, ShellEvent::Close);
                    if events_tx.send(event).await.is_err() || close {
                        break;
                    }
                }
            });

            tokio::spawn(async move {
                while let Some(command) = commands_rx.recv().await {
                    let result = match command {
                        ShellCommand::Data(data) => write_half.data(Cursor::new(data)).await,
                        ShellCommand::Resize(size) => {
                            write_half
                                .window_change(
                                    size.col_width,
                                    size.row_height,
                                    size.pix_width,
                                    size.pix_height,
                                )
                                .await
                        }
                        ShellCommand::Eof => write_half.eof().await,
                        ShellCommand::Close => {
                            let result = write_half.close().await;
                            if result.is_err() {
                                break;
                            }
                            break;
                        }
                    };
                    if result.is_err() {
                        break;
                    }
                }
            });

            Ok(OpenedShell {
                session: Box::new(GuestShellSession {
                    _session: session,
                    commands: commands_tx,
                }),
                events: events_rx,
            })
        })
    }
}

async fn connect_guest(cfg: &Config, addr: SocketAddr) -> Result<client::Handle<GuestClient>> {
    let key = load_secret_key(&cfg.guest_key, None)
        .with_context(|| format!("load guest SSH key {}", cfg.guest_key.display()))?;
    let mut session = client::connect(Arc::new(client::Config::default()), addr, GuestClient)
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
    Ok(session)
}

enum ShellCommand {
    Data(Vec<u8>),
    Resize(TerminalSize),
    Eof,
    Close,
}

struct GuestShellSession {
    _session: client::Handle<GuestClient>,
    commands: mpsc::Sender<ShellCommand>,
}

impl ShellSession for GuestShellSession {
    fn data(&mut self, data: Vec<u8>) -> BoxFuture<Result<()>> {
        let commands = self.commands.clone();
        Box::pin(async move {
            commands
                .send(ShellCommand::Data(data))
                .await
                .context("send data to guest shell task")
        })
    }

    fn resize(&mut self, size: TerminalSize) -> BoxFuture<Result<()>> {
        let commands = self.commands.clone();
        Box::pin(async move {
            commands
                .send(ShellCommand::Resize(size))
                .await
                .context("send window-change to guest shell task")
        })
    }

    fn eof(&mut self) -> BoxFuture<Result<()>> {
        let commands = self.commands.clone();
        Box::pin(async move {
            commands
                .send(ShellCommand::Eof)
                .await
                .context("send EOF to guest shell task")
        })
    }

    fn close(&mut self) -> BoxFuture<Result<()>> {
        let commands = self.commands.clone();
        Box::pin(async move {
            commands
                .send(ShellCommand::Close)
                .await
                .context("send close to guest shell task")
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
            .exec(b"echo hello".to_vec(), None)
            .await
            .expect("proxy exec");

        assert_eq!(output.stdout, b"ran: echo hello\n");
        assert_eq!(output.stderr, b"warn\n");
        assert_eq!(output.exit_status, 7);
        server_task.await.expect("server task");
        fs::remove_dir_all(key_path.parent().unwrap()).expect("remove key dir");
    }

    #[tokio::test]
    async fn guest_proxy_requests_pty_before_exec_when_provided() {
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
        let pty = PtyRequest::new(
            "xterm-256color".to_string(),
            TerminalSize::new(120, 40, 800, 600),
            vec![(Pty::ECHO, 1)],
        );

        let output = proxy
            .exec(b"echo hello".to_vec(), Some(pty))
            .await
            .expect("proxy exec");

        assert_eq!(
            output.stdout,
            b"pty xterm-256color 120 40 800 600\nran: echo hello\n"
        );
        assert_eq!(output.stderr, b"warn\n");
        assert_eq!(output.exit_status, 7);
        server_task.await.expect("server task");
        fs::remove_dir_all(key_path.parent().unwrap()).expect("remove key dir");
    }

    #[tokio::test]
    async fn guest_proxy_shell_forwards_pty_data_window_and_eof() {
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
        let pty = PtyRequest::new(
            "xterm-256color".to_string(),
            TerminalSize::new(80, 24, 640, 480),
            vec![(Pty::ECHO, 1)],
        );

        let OpenedShell {
            mut session,
            mut events,
        } = proxy.open_shell(Some(pty)).await.expect("open shell");
        session
            .resize(TerminalSize::new(120, 40, 960, 640))
            .await
            .expect("resize shell");
        session
            .data(b"hello shell\n".to_vec())
            .await
            .expect("send shell data");
        session.eof().await.expect("send shell eof");

        let mut stdout = Vec::new();
        let mut exit_status = None;
        while let Some(event) = events.recv().await {
            match event {
                ShellEvent::Stdout(data) => stdout.extend_from_slice(&data),
                ShellEvent::ExitStatus(status) => exit_status = Some(status),
                ShellEvent::Close => break,
                _ => {}
            }
        }

        assert_eq!(
            stdout,
            b"shell xterm-256color 120 40 960 640\ninput: hello shell\n"
        );
        assert_eq!(exit_status, Some(0));
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
                pty: None,
            }
        }
    }

    struct InnerSession {
        accepted: PublicKey,
        pty: Option<PtyRequest>,
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

        async fn pty_request(
            &mut self,
            channel: ChannelId,
            term: &str,
            col_width: u32,
            row_height: u32,
            pix_width: u32,
            pix_height: u32,
            modes: &[(Pty, u32)],
            session: &mut Session,
        ) -> Result<(), Self::Error> {
            self.pty = Some(PtyRequest::new(
                term.to_string(),
                TerminalSize::new(col_width, row_height, pix_width, pix_height),
                modes.to_vec(),
            ));
            session.channel_success(channel)?;
            Ok(())
        }

        async fn exec_request(
            &mut self,
            channel: ChannelId,
            data: &[u8],
            session: &mut Session,
        ) -> Result<(), Self::Error> {
            session.channel_success(channel)?;
            if let Some(pty) = &self.pty {
                session.data(
                    channel,
                    format!(
                        "pty {} {} {} {} {}\n",
                        pty.term,
                        pty.size.col_width,
                        pty.size.row_height,
                        pty.size.pix_width,
                        pty.size.pix_height
                    ),
                )?;
            }
            session.data(channel, format!("ran: {}\n", String::from_utf8_lossy(data)))?;
            session.extended_data(channel, 1, "warn\n")?;
            session.exit_status_request(channel, 7)?;
            session.eof(channel)?;
            session.close(channel)?;
            Ok(())
        }

        async fn shell_request(
            &mut self,
            channel: ChannelId,
            session: &mut Session,
        ) -> Result<(), Self::Error> {
            session.channel_success(channel)?;
            Ok(())
        }

        async fn window_change_request(
            &mut self,
            _channel: ChannelId,
            col_width: u32,
            row_height: u32,
            pix_width: u32,
            pix_height: u32,
            _session: &mut Session,
        ) -> Result<(), Self::Error> {
            if let Some(pty) = self.pty.as_mut() {
                pty.resize(TerminalSize::new(
                    col_width, row_height, pix_width, pix_height,
                ));
            }
            Ok(())
        }

        async fn data(
            &mut self,
            channel: ChannelId,
            data: &[u8],
            session: &mut Session,
        ) -> Result<(), Self::Error> {
            if let Some(pty) = &self.pty {
                session.data(
                    channel,
                    format!(
                        "shell {} {} {} {} {}\n",
                        pty.term,
                        pty.size.col_width,
                        pty.size.row_height,
                        pty.size.pix_width,
                        pty.size.pix_height
                    ),
                )?;
            }
            session.data(channel, format!("input: {}", String::from_utf8_lossy(data)))?;
            Ok(())
        }

        async fn channel_eof(
            &mut self,
            channel: ChannelId,
            session: &mut Session,
        ) -> Result<(), Self::Error> {
            session.exit_status_request(channel, 0)?;
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
