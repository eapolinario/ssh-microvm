use std::{
    env,
    ffi::OsString,
    fs,
    net::{SocketAddr, TcpListener, TcpStream},
    path::PathBuf,
    process::{Child, Command, Output, Stdio},
    thread,
    time::{Duration, Instant, SystemTime, UNIX_EPOCH},
};

const SERVER_START_TIMEOUT: Duration = Duration::from_secs(5);
const SSH_COMMAND_TIMEOUT: Duration = Duration::from_secs(120);

#[test]
fn boots_vm_and_execs_echo_hello() {
    let Some(env) = SmokeEnv::from_env() else {
        eprintln!(
            "skipping smoke test; set SSH_MICROVM_KERNEL, SSH_MICROVM_ROOTFS, and SSH_MICROVM_GUEST_KEY"
        );
        return;
    };

    let listen = unused_local_addr();
    let state_dir = unique_temp_dir();
    fs::create_dir_all(&state_dir).expect("create smoke test state directory");

    let mut server = ServerGuard::spawn(&env, listen, state_dir.clone());
    wait_for_server(listen, &mut server);

    let output = run_ssh_command(&env.guest_key, listen).expect("run ssh smoke command");
    assert!(
        output.status.success(),
        "ssh command should exit successfully\nstdout:\n{}\nstderr:\n{}",
        String::from_utf8_lossy(&output.stdout),
        String::from_utf8_lossy(&output.stderr)
    );
    assert_eq!(String::from_utf8_lossy(&output.stdout).trim(), "hello");

    server.stop();
    let _ = fs::remove_dir_all(state_dir);
}

struct SmokeEnv {
    kernel: OsString,
    rootfs: OsString,
    guest_key: OsString,
}

impl SmokeEnv {
    fn from_env() -> Option<Self> {
        Some(Self {
            kernel: env::var_os("SSH_MICROVM_KERNEL")?,
            rootfs: env::var_os("SSH_MICROVM_ROOTFS")?,
            guest_key: env::var_os("SSH_MICROVM_GUEST_KEY")?,
        })
    }
}

struct ServerGuard {
    child: Option<Child>,
}

impl ServerGuard {
    fn spawn(env: &SmokeEnv, listen: SocketAddr, state_dir: PathBuf) -> Self {
        let child = Command::new(env!("CARGO_BIN_EXE_ssh-microvm"))
            .arg("--listen")
            .arg(listen.to_string())
            .arg("--kernel")
            .arg(&env.kernel)
            .arg("--rootfs")
            .arg(&env.rootfs)
            .arg("--guest-key")
            .arg(&env.guest_key)
            .arg("--state-dir")
            .arg(state_dir)
            .arg("--accept-any-key")
            .arg("--boot-timeout")
            .arg("60s")
            .stdin(Stdio::null())
            .stdout(Stdio::piped())
            .stderr(Stdio::piped())
            .spawn()
            .expect("spawn ssh-microvm server");

        Self { child: Some(child) }
    }

    fn try_wait(&mut self) -> Option<std::process::ExitStatus> {
        self.child
            .as_mut()
            .and_then(|child| child.try_wait().expect("poll ssh-microvm server"))
    }

    fn stop(&mut self) {
        let Some(mut child) = self.child.take() else {
            return;
        };

        if child.try_wait().expect("poll ssh-microvm server").is_none() {
            child.kill().expect("stop ssh-microvm server");
        }
        let _ = child.wait();
    }
}

impl Drop for ServerGuard {
    fn drop(&mut self) {
        self.stop();
    }
}

fn wait_for_server(listen: SocketAddr, server: &mut ServerGuard) {
    let deadline = Instant::now() + SERVER_START_TIMEOUT;
    loop {
        if TcpStream::connect(listen).is_ok() {
            return;
        }

        if let Some(status) = server.try_wait() {
            panic!("ssh-microvm server exited before accepting connections: {status}");
        }

        assert!(
            Instant::now() < deadline,
            "ssh-microvm server did not start listening on {listen} within {SERVER_START_TIMEOUT:?}"
        );
        thread::sleep(Duration::from_millis(50));
    }
}

fn run_ssh_command(guest_key: &OsString, listen: SocketAddr) -> Result<Output, String> {
    let mut child = Command::new("ssh")
        .arg("-F")
        .arg("/dev/null")
        .arg("-o")
        .arg("BatchMode=yes")
        .arg("-o")
        .arg("IdentitiesOnly=yes")
        .arg("-o")
        .arg("LogLevel=ERROR")
        .arg("-o")
        .arg("StrictHostKeyChecking=no")
        .arg("-o")
        .arg("UserKnownHostsFile=/dev/null")
        .arg("-i")
        .arg(guest_key)
        .arg("-p")
        .arg(listen.port().to_string())
        .arg(format!("root@{}", listen.ip()))
        .arg("echo hello")
        .stdin(Stdio::null())
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .spawn()
        .map_err(|err| format!("spawn ssh client: {err}"))?;

    let deadline = Instant::now() + SSH_COMMAND_TIMEOUT;
    loop {
        if child
            .try_wait()
            .map_err(|err| format!("poll ssh client: {err}"))?
            .is_some()
        {
            return child
                .wait_with_output()
                .map_err(|err| format!("collect ssh client output: {err}"));
        }

        if Instant::now() >= deadline {
            child
                .kill()
                .map_err(|err| format!("kill timed-out ssh client: {err}"))?;
            let output = child
                .wait_with_output()
                .map_err(|err| format!("collect timed-out ssh client output: {err}"))?;
            return Err(format!(
                "ssh command timed out after {SSH_COMMAND_TIMEOUT:?}\nstdout:\n{}\nstderr:\n{}",
                String::from_utf8_lossy(&output.stdout),
                String::from_utf8_lossy(&output.stderr)
            ));
        }

        thread::sleep(Duration::from_millis(100));
    }
}

fn unused_local_addr() -> SocketAddr {
    let listener = TcpListener::bind("127.0.0.1:0").expect("reserve local TCP port");
    listener.local_addr().expect("read reserved local TCP port")
}

fn unique_temp_dir() -> PathBuf {
    let mut dir = env::temp_dir();
    let nanos = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .expect("system clock should be after Unix epoch")
        .as_nanos();
    dir.push(format!("ssh-microvm-smoke-{}-{nanos}", std::process::id()));
    dir
}
