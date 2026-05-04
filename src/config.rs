//! Command-line configuration for ssh-microvm.

use std::{
    net::{IpAddr, SocketAddr},
    path::PathBuf,
    time::Duration,
};

use clap::{ArgGroup, Parser, Subcommand};

#[derive(Debug, Parser)]
#[command(name = "ssh-microvm", version, about)]
#[command(group(
    ArgGroup::new("auth")
        .required(true)
        .multiple(false)
        .args(["authorized_keys", "accept_any_key"])
))]
pub struct Config {
    /// Boot one Firecracker microVM, wait for guest sshd, then tear it down.
    #[arg(long)]
    pub dry_boot: bool,

    #[command(subcommand)]
    pub command: Option<Command>,

    /// SSH listen address for incoming client connections.
    #[arg(long, default_value = "0.0.0.0:2222")]
    pub listen: SocketAddr,

    /// Path to the guest Linux kernel image.
    #[arg(long, value_name = "PATH")]
    pub kernel: PathBuf,

    /// Path to the guest ext4 root filesystem.
    #[arg(long, value_name = "PATH")]
    pub rootfs: PathBuf,

    /// Directory for per-VM sockets, logs, and generated host keys.
    #[arg(long, value_name = "DIR", default_value = "/tmp/ssh-microvm")]
    pub state_dir: PathBuf,

    /// SSH host key path. Defaults to <state-dir>/ssh_host_ed25519.
    #[arg(long, value_name = "PATH")]
    pub host_key: Option<PathBuf>,

    /// OpenSSH authorized_keys file used for client public-key authentication.
    #[arg(long, value_name = "PATH")]
    pub authorized_keys: Option<PathBuf>,

    /// Accept any client public key. Intended for local development only.
    #[arg(long)]
    pub accept_any_key: bool,

    /// Path to the Firecracker executable.
    #[arg(long, value_name = "PATH", default_value = "firecracker")]
    pub firecracker: PathBuf,

    /// Number of virtual CPUs assigned to the guest.
    #[arg(long, default_value_t = 1)]
    pub vcpu: u8,

    /// Guest memory size in MiB.
    #[arg(long = "mem", default_value_t = 512)]
    pub mem_mib: u32,

    /// Kernel command line passed to the guest.
    #[arg(long, default_value = "console=ttyS0 reboot=k panic=1 pci=off")]
    pub boot_args: String,

    /// SSH username used when connecting to guest sshd.
    #[arg(long, default_value = "root")]
    pub guest_user: String,

    /// Private key used to connect to guest sshd.
    #[arg(long, value_name = "PATH")]
    pub guest_key: PathBuf,

    /// Guest IP address reachable over the host tap device.
    #[arg(long, default_value = "172.16.0.2")]
    pub guest_ip: IpAddr,

    /// Host IP address configured on the tap device.
    #[arg(long, default_value = "172.16.0.1")]
    pub host_ip: IpAddr,

    /// Host tap device name. The device is created externally.
    #[arg(long, default_value = "tap0")]
    pub tap_name: String,

    /// How long to wait for guest sshd to become reachable.
    #[arg(long, value_parser = parse_duration, default_value = "30s")]
    pub boot_timeout: Duration,

    /// Grace period after requesting shutdown before force-killing Firecracker.
    #[arg(long, value_parser = parse_duration, default_value = "2s")]
    pub grace_stop: Duration,
}

#[derive(Debug, Clone, Subcommand, PartialEq, Eq)]
pub enum Command {
    /// MicroVM lifecycle commands.
    Microvm {
        #[command(subcommand)]
        command: MicrovmCommand,
    },
}

#[derive(Debug, Clone, Subcommand, PartialEq, Eq)]
pub enum MicrovmCommand {
    /// Boot one Firecracker microVM and wait for guest sshd.
    Boot,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum RunMode {
    Server,
    BootOnce,
}

impl Config {
    pub fn host_key_path(&self) -> PathBuf {
        self.host_key
            .clone()
            .unwrap_or_else(|| self.state_dir.join("ssh_host_ed25519"))
    }

    pub fn run_mode(&self) -> RunMode {
        if self.dry_boot
            || matches!(
                self.command,
                Some(Command::Microvm {
                    command: MicrovmCommand::Boot,
                })
            )
        {
            RunMode::BootOnce
        } else {
            RunMode::Server
        }
    }
}

fn parse_duration(raw: &str) -> Result<Duration, String> {
    let raw = raw.trim();
    if raw.is_empty() {
        return Err("duration must not be empty".to_string());
    }

    if let Some(value) = raw.strip_suffix("ms") {
        return parse_duration_number(value).map(Duration::from_millis);
    }

    if let Some(value) = raw.strip_suffix('s') {
        return parse_duration_number(value).map(Duration::from_secs);
    }

    if let Some(value) = raw.strip_suffix('m') {
        return parse_duration_number(value).map(|minutes| Duration::from_secs(minutes * 60));
    }

    parse_duration_number(raw).map(Duration::from_secs)
}

fn parse_duration_number(raw: &str) -> Result<u64, String> {
    raw.parse::<u64>()
        .map_err(|err| format!("invalid duration value {raw:?}: {err}"))
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::ffi::OsString;

    use clap::{CommandFactory, error::ErrorKind};

    fn base_args() -> Vec<&'static str> {
        vec![
            "ssh-microvm",
            "--kernel",
            "artifacts/vmlinux.bin",
            "--rootfs",
            "artifacts/ubuntu.ext4",
            "--guest-key",
            "artifacts/ubuntu.id_rsa",
            "--accept-any-key",
        ]
    }

    #[test]
    fn parses_required_paths_and_defaults() {
        let cfg = Config::try_parse_from(base_args()).expect("valid config");

        assert_eq!(cfg.listen, "0.0.0.0:2222".parse().unwrap());
        assert_eq!(cfg.run_mode(), RunMode::Server);
        assert_eq!(cfg.kernel, PathBuf::from("artifacts/vmlinux.bin"));
        assert_eq!(cfg.rootfs, PathBuf::from("artifacts/ubuntu.ext4"));
        assert_eq!(cfg.state_dir, PathBuf::from("/tmp/ssh-microvm"));
        assert_eq!(
            cfg.host_key_path(),
            PathBuf::from("/tmp/ssh-microvm/ssh_host_ed25519")
        );
        assert!(cfg.accept_any_key);
        assert_eq!(cfg.firecracker, PathBuf::from("firecracker"));
        assert_eq!(cfg.vcpu, 1);
        assert_eq!(cfg.mem_mib, 512);
        assert_eq!(cfg.boot_args, "console=ttyS0 reboot=k panic=1 pci=off");
        assert_eq!(cfg.guest_user, "root");
        assert_eq!(cfg.guest_key, PathBuf::from("artifacts/ubuntu.id_rsa"));
        assert_eq!(cfg.guest_ip, "172.16.0.2".parse::<IpAddr>().unwrap());
        assert_eq!(cfg.host_ip, "172.16.0.1".parse::<IpAddr>().unwrap());
        assert_eq!(cfg.tap_name, "tap0");
        assert_eq!(cfg.boot_timeout, Duration::from_secs(30));
        assert_eq!(cfg.grace_stop, Duration::from_secs(2));
    }

    #[test]
    fn parses_overrides() {
        let mut args = base_args();
        args.extend([
            "--listen",
            "127.0.0.1:2223",
            "--state-dir",
            "/var/run/ssh-microvm",
            "--host-key",
            "/keys/host",
            "--firecracker",
            "/usr/bin/firecracker",
            "--vcpu",
            "2",
            "--mem",
            "1024",
            "--boot-args",
            "console=ttyS0",
            "--guest-user",
            "ubuntu",
            "--guest-ip",
            "10.0.0.2",
            "--host-ip",
            "10.0.0.1",
            "--tap-name",
            "tap1",
            "--boot-timeout",
            "45s",
            "--grace-stop",
            "1500ms",
        ]);

        let cfg = Config::try_parse_from(args).expect("valid config with overrides");

        assert_eq!(cfg.listen, "127.0.0.1:2223".parse().unwrap());
        assert_eq!(cfg.state_dir, PathBuf::from("/var/run/ssh-microvm"));
        assert_eq!(cfg.host_key_path(), PathBuf::from("/keys/host"));
        assert_eq!(cfg.firecracker, PathBuf::from("/usr/bin/firecracker"));
        assert_eq!(cfg.vcpu, 2);
        assert_eq!(cfg.mem_mib, 1024);
        assert_eq!(cfg.boot_args, "console=ttyS0");
        assert_eq!(cfg.guest_user, "ubuntu");
        assert_eq!(cfg.guest_ip, "10.0.0.2".parse::<IpAddr>().unwrap());
        assert_eq!(cfg.host_ip, "10.0.0.1".parse::<IpAddr>().unwrap());
        assert_eq!(cfg.tap_name, "tap1");
        assert_eq!(cfg.boot_timeout, Duration::from_secs(45));
        assert_eq!(cfg.grace_stop, Duration::from_millis(1500));
    }

    #[test]
    fn requires_one_auth_mode() {
        let args = [
            "ssh-microvm",
            "--kernel",
            "artifacts/vmlinux.bin",
            "--rootfs",
            "artifacts/ubuntu.ext4",
            "--guest-key",
            "artifacts/ubuntu.id_rsa",
        ];

        let err = Config::try_parse_from(args).expect_err("auth mode should be required");

        assert_eq!(err.kind(), ErrorKind::MissingRequiredArgument);
    }

    #[test]
    fn rejects_multiple_auth_modes() {
        let mut args = base_args();
        args.extend(["--authorized-keys", "authorized_keys"]);

        let err = Config::try_parse_from(args).expect_err("auth modes should conflict");

        assert_eq!(err.kind(), ErrorKind::ArgumentConflict);
    }

    #[test]
    fn help_documents_every_flag() {
        let mut help = Vec::new();
        Config::command()
            .write_long_help(&mut help)
            .expect("render help");
        let help = String::from_utf8(help).expect("help is valid utf-8");

        for flag in [
            "--listen",
            "--dry-boot",
            "--kernel",
            "--rootfs",
            "--state-dir",
            "--host-key",
            "--authorized-keys",
            "--accept-any-key",
            "--firecracker",
            "--vcpu",
            "--mem",
            "--boot-args",
            "--guest-user",
            "--guest-key",
            "--guest-ip",
            "--host-ip",
            "--tap-name",
            "--boot-timeout",
            "--grace-stop",
        ] {
            assert!(help.contains(flag), "help should document {flag}");
        }
    }

    #[test]
    fn dry_boot_flag_selects_boot_once_mode() {
        let mut args = base_args();
        args.push("--dry-boot");

        let cfg = Config::try_parse_from(args).expect("valid dry boot config");

        assert_eq!(cfg.run_mode(), RunMode::BootOnce);
    }

    #[test]
    fn microvm_boot_subcommand_selects_boot_once_mode() {
        let mut args = base_args();
        args.extend(["microvm", "boot"]);

        let cfg = Config::try_parse_from(args).expect("valid microvm boot config");

        assert_eq!(cfg.run_mode(), RunMode::BootOnce);
    }

    #[test]
    fn duration_parser_rejects_empty_values() {
        let err = parse_duration("  ").expect_err("empty duration should fail");

        assert_eq!(err, "duration must not be empty");
    }

    #[test]
    fn duration_parser_rejects_unknown_units() {
        let err = parse_duration("1h").expect_err("unsupported duration unit should fail");

        assert!(err.contains("invalid duration value"));
    }

    #[test]
    fn accepts_os_strings_from_clap() {
        let args = base_args()
            .into_iter()
            .map(OsString::from)
            .collect::<Vec<_>>();

        Config::try_parse_from(args).expect("os string args should parse");
    }
}
