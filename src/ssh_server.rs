//! Outer SSH server.

use std::{path::Path, sync::Arc};

use anyhow::{Context, Result, anyhow};
use russh::{
    Channel, ChannelId, MethodKind, MethodSet,
    keys::{
        Algorithm, PrivateKey, PublicKey,
        key::safe_rng,
        load_secret_key,
        ssh_key::{AuthorizedKeys, LineEnding},
    },
    server::{self, Auth, Msg, Server as _, Session},
};
use tokio::net::TcpListener;

use crate::{
    config::Config,
    lifecycle::{VmLease, VmLifecycle},
    proxy::{ExecProxy, GuestSshProxy},
};

#[derive(Clone, Debug)]
enum ClientAuth {
    AcceptAnyKey,
    AuthorizedKeys(Arc<Vec<PublicKey>>),
}

impl ClientAuth {
    fn accepts(&self, public_key: &PublicKey) -> bool {
        match self {
            Self::AcceptAnyKey => true,
            Self::AuthorizedKeys(keys) => keys.iter().any(|key| key == public_key),
        }
    }
}

struct SshServer {
    auth: ClientAuth,
    lifecycle: VmLifecycle,
    proxy: Arc<dyn ExecProxy>,
}

impl server::Server for SshServer {
    type Handler = SshSession;

    fn new_client(&mut self, peer_addr: Option<std::net::SocketAddr>) -> Self::Handler {
        tracing::debug!(?peer_addr, "accepted ssh client");
        SshSession {
            auth: self.auth.clone(),
            lifecycle: self.lifecycle.clone(),
            proxy: self.proxy.clone(),
            vm_lease: None,
        }
    }

    fn handle_session_error(&mut self, error: <Self::Handler as server::Handler>::Error) {
        tracing::debug!(error = ?error, "ssh session ended with error");
    }
}

struct SshSession {
    auth: ClientAuth,
    lifecycle: VmLifecycle,
    proxy: Arc<dyn ExecProxy>,
    vm_lease: Option<VmLease>,
}

impl server::Handler for SshSession {
    type Error = anyhow::Error;

    async fn auth_publickey_offered(
        &mut self,
        _user: &str,
        public_key: &PublicKey,
    ) -> Result<Auth, Self::Error> {
        Ok(auth_result(self.auth.accepts(public_key)))
    }

    async fn auth_publickey(
        &mut self,
        _user: &str,
        public_key: &PublicKey,
    ) -> Result<Auth, Self::Error> {
        Ok(auth_result(self.auth.accepts(public_key)))
    }

    async fn channel_open_session(
        &mut self,
        _channel: Channel<Msg>,
        _session: &mut Session,
    ) -> Result<bool, Self::Error> {
        if self.vm_lease.is_none() {
            self.vm_lease = Some(
                self.lifecycle
                    .acquire()
                    .await
                    .context("boot VM for SSH session channel")?,
            );
        }
        Ok(true)
    }

    async fn exec_request(
        &mut self,
        channel: ChannelId,
        data: &[u8],
        session: &mut Session,
    ) -> Result<(), Self::Error> {
        let output = match self.proxy.exec(data.to_vec()).await {
            Ok(output) => output,
            Err(err) => {
                session.channel_failure(channel)?;
                session.close(channel)?;
                return Err(err.context("proxy guest SSH exec"));
            }
        };

        session.channel_success(channel)?;
        if !output.stdout.is_empty() {
            session.data(channel, output.stdout)?;
        }
        if !output.stderr.is_empty() {
            session.extended_data(channel, 1, output.stderr)?;
        }
        session.exit_status_request(channel, output.exit_status)?;
        session.eof(channel)?;
        session.close(channel)?;
        Ok(())
    }
}

fn auth_result(accepted: bool) -> Auth {
    if accepted {
        Auth::Accept
    } else {
        Auth::reject()
    }
}

pub async fn run(cfg: &Config) -> Result<()> {
    let host_key = load_or_generate_host_key(&cfg.host_key_path())?;
    let auth = load_client_auth(cfg)?;
    let lifecycle = VmLifecycle::new(cfg);
    let proxy = Arc::new(GuestSshProxy::new(cfg));
    let server_config = Arc::new(server_config(host_key));
    let listener = TcpListener::bind(cfg.listen)
        .await
        .with_context(|| format!("bind SSH listener on {}", cfg.listen))?;
    let mut server = SshServer {
        auth,
        lifecycle,
        proxy,
    };

    tracing::info!(listen = %cfg.listen, "ssh server listening");
    server.run_on_socket(server_config, &listener).await?;
    Ok(())
}

fn server_config(host_key: PrivateKey) -> server::Config {
    server::Config {
        keys: vec![host_key],
        methods: MethodSet::from(&[MethodKind::PublicKey][..]),
        auth_rejection_time_initial: Some(std::time::Duration::from_millis(0)),
        ..Default::default()
    }
}

fn load_client_auth(cfg: &Config) -> Result<ClientAuth> {
    if cfg.accept_any_key {
        return Ok(ClientAuth::AcceptAnyKey);
    }

    let path = cfg
        .authorized_keys
        .as_ref()
        .ok_or_else(|| anyhow!("--authorized-keys is required unless --accept-any-key is set"))?;
    let contents = std::fs::read_to_string(path)
        .with_context(|| format!("read authorized keys {}", path.display()))?;
    let keys = AuthorizedKeys::new(&contents)
        .map(|entry| entry.map(|entry| entry.public_key().clone()))
        .collect::<std::result::Result<Vec<_>, _>>()
        .with_context(|| format!("parse authorized keys {}", path.display()))?;

    Ok(ClientAuth::AuthorizedKeys(Arc::new(keys)))
}

fn load_or_generate_host_key(path: &Path) -> Result<PrivateKey> {
    if path.exists() {
        return load_secret_key(path, None)
            .with_context(|| format!("load SSH host key {}", path.display()));
    }

    if let Some(parent) = path.parent() {
        std::fs::create_dir_all(parent)
            .with_context(|| format!("create host key directory {}", parent.display()))?;
    }

    let key =
        PrivateKey::random(&mut safe_rng(), Algorithm::Ed25519).context("generate SSH host key")?;
    write_private_key(path, &key)?;
    Ok(key)
}

fn write_private_key(path: &Path, key: &PrivateKey) -> Result<()> {
    let encoded = key
        .to_openssh(LineEnding::LF)
        .context("encode SSH host key")?;

    #[cfg(unix)]
    {
        use std::io::Write as _;
        use std::os::unix::fs::OpenOptionsExt as _;

        let mut file = std::fs::OpenOptions::new()
            .write(true)
            .create_new(true)
            .mode(0o600)
            .open(path)
            .with_context(|| format!("create SSH host key {}", path.display()))?;
        file.write_all(encoded.as_bytes())
            .with_context(|| format!("write SSH host key {}", path.display()))?;
    }

    #[cfg(not(unix))]
    {
        std::fs::write(path, encoded.as_bytes())
            .with_context(|| format!("write SSH host key {}", path.display()))?;
    }

    Ok(())
}

#[cfg(test)]
mod tests {
    use std::{
        fs,
        net::{IpAddr, Ipv4Addr, SocketAddr},
        path::PathBuf,
        sync::{
            Arc,
            atomic::{AtomicUsize, Ordering},
        },
        time::Duration,
    };

    use russh::{
        ChannelMsg, Disconnect,
        client::{self, Handler},
        keys::PrivateKeyWithHashAlg,
    };
    use tokio::{sync::Notify, time::timeout};

    use crate::{
        lifecycle::{BoxFuture, RunningVm, VmLifecycle, VmSpawner},
        proxy::ExecOutput,
    };

    use super::*;

    fn test_config() -> Config {
        Config {
            dry_boot: false,
            command: None,
            listen: SocketAddr::new(IpAddr::V4(Ipv4Addr::LOCALHOST), 0),
            kernel: PathBuf::from("kernel"),
            rootfs: PathBuf::from("rootfs"),
            state_dir: unique_test_dir("state"),
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

    fn unique_test_dir(name: &str) -> PathBuf {
        std::env::temp_dir().join(format!(
            "ssh-microvm-{name}-{}-{}",
            std::process::id(),
            std::time::SystemTime::now()
                .duration_since(std::time::UNIX_EPOCH)
                .expect("system time should be after epoch")
                .as_nanos()
        ))
    }

    #[test]
    fn generates_and_reloads_host_key() {
        let dir = unique_test_dir("host-key");
        let path = dir.join("ssh_host_ed25519");

        let generated = load_or_generate_host_key(&path).expect("generate host key");
        let reloaded = load_or_generate_host_key(&path).expect("reload host key");

        assert_eq!(generated.public_key(), reloaded.public_key());

        #[cfg(unix)]
        {
            use std::os::unix::fs::PermissionsExt as _;
            let mode = fs::metadata(&path)
                .expect("host key metadata")
                .permissions()
                .mode()
                & 0o777;
            assert_eq!(mode, 0o600);
        }

        fs::remove_dir_all(dir).expect("remove test dir");
    }

    #[test]
    fn authorized_keys_auth_accepts_only_listed_keys() {
        let accepted =
            PrivateKey::random(&mut safe_rng(), Algorithm::Ed25519).expect("generate accepted key");
        let rejected =
            PrivateKey::random(&mut safe_rng(), Algorithm::Ed25519).expect("generate rejected key");
        let path = unique_test_dir("authorized-keys").join("authorized_keys");
        fs::create_dir_all(path.parent().expect("path has parent")).expect("create dir");
        fs::write(
            &path,
            format!(
                "{}\n",
                accepted
                    .public_key()
                    .to_openssh()
                    .expect("encode public key")
            ),
        )
        .expect("write authorized_keys");

        let mut cfg = test_config();
        cfg.accept_any_key = false;
        cfg.authorized_keys = Some(path.clone());
        let auth = load_client_auth(&cfg).expect("load authorized_keys");

        assert!(auth.accepts(accepted.public_key()));
        assert!(!auth.accepts(rejected.public_key()));

        fs::remove_dir_all(path.parent().unwrap()).expect("remove test dir");
    }

    #[tokio::test]
    async fn ssh_exec_replies_ok() {
        let client_key =
            PrivateKey::random(&mut safe_rng(), Algorithm::Ed25519).expect("generate client key");
        let host_key =
            PrivateKey::random(&mut safe_rng(), Algorithm::Ed25519).expect("generate host key");
        let listener = TcpListener::bind("127.0.0.1:0")
            .await
            .expect("bind test server");
        let addr = listener.local_addr().expect("listener addr");
        let server_config = Arc::new(server_config(host_key));
        let cfg = test_config();
        let vm_spawner = Arc::new(TestVmSpawner::default());
        let lifecycle = VmLifecycle::new_with_spawner(&cfg, vm_spawner.clone());
        let mut server = SshServer {
            auth: ClientAuth::AcceptAnyKey,
            lifecycle,
            proxy: Arc::new(TestProxy),
        };

        let server_task = tokio::spawn(async move {
            let (socket, peer_addr) = listener.accept().await.expect("accept client");
            let handler = server.new_client(Some(peer_addr));
            server::run_stream(server_config, socket, handler)
                .await
                .expect("run ssh session");
        });

        let mut session = client::connect(Arc::new(client::Config::default()), addr, TestClient)
            .await
            .expect("connect client");
        let authenticated = session
            .authenticate_publickey(
                "test",
                PrivateKeyWithHashAlg::new(
                    Arc::new(client_key),
                    session
                        .best_supported_rsa_hash()
                        .await
                        .expect("query rsa hash")
                        .flatten(),
                ),
            )
            .await
            .expect("authenticate")
            .success();
        assert!(authenticated);

        let mut channel = session
            .channel_open_session()
            .await
            .expect("open session channel");
        channel
            .exec(true, "echo through guest")
            .await
            .expect("exec request");

        let mut stdout = Vec::new();
        let mut exit_status = None;
        while let Some(msg) = channel.wait().await {
            match msg {
                ChannelMsg::Data { data } => stdout.extend_from_slice(&data),
                ChannelMsg::ExitStatus { exit_status: code } => exit_status = Some(code),
                ChannelMsg::Close => break,
                _ => {}
            }
        }

        assert_eq!(stdout, b"proxied: echo through guest\n");
        assert_eq!(exit_status, Some(3));
        assert_eq!(vm_spawner.boots.load(Ordering::SeqCst), 1);

        session
            .disconnect(Disconnect::ByApplication, "", "en")
            .await
            .expect("disconnect");
        server_task.await.expect("server task");
        timeout(
            Duration::from_secs(1),
            vm_spawner.shutdown_notify.notified(),
        )
        .await
        .expect("disconnect should shut down the VM lease");
        assert_eq!(vm_spawner.shutdowns.load(Ordering::SeqCst), 1);
    }

    struct TestClient;

    impl Handler for TestClient {
        type Error = anyhow::Error;

        async fn check_server_key(&mut self, _server_public_key: &PublicKey) -> Result<bool> {
            Ok(true)
        }
    }

    struct TestProxy;

    impl ExecProxy for TestProxy {
        fn exec(&self, command: Vec<u8>) -> BoxFuture<Result<ExecOutput>> {
            Box::pin(async move {
                Ok(ExecOutput {
                    stdout: format!("proxied: {}\n", String::from_utf8_lossy(&command))
                        .into_bytes(),
                    stderr: Vec::new(),
                    exit_status: 3,
                })
            })
        }
    }

    #[derive(Default)]
    struct TestVmSpawner {
        boots: AtomicUsize,
        shutdowns: Arc<AtomicUsize>,
        shutdown_notify: Arc<Notify>,
    }

    impl VmSpawner for TestVmSpawner {
        fn boot(&self, _cfg: Config) -> BoxFuture<Result<Box<dyn RunningVm>>> {
            self.boots.fetch_add(1, Ordering::SeqCst);
            let shutdowns = self.shutdowns.clone();
            let shutdown_notify = self.shutdown_notify.clone();
            Box::pin(async move {
                Ok(Box::new(TestVm {
                    shutdowns,
                    shutdown_notify,
                }) as Box<dyn RunningVm>)
            })
        }
    }

    struct TestVm {
        shutdowns: Arc<AtomicUsize>,
        shutdown_notify: Arc<Notify>,
    }

    impl RunningVm for TestVm {
        fn shutdown(self: Box<Self>) -> BoxFuture<Result<()>> {
            Box::pin(async move {
                self.shutdowns.fetch_add(1, Ordering::SeqCst);
                self.shutdown_notify.notify_one();
                Ok(())
            })
        }
    }
}
