use std::process::Command;

#[test]
fn help_lists_config_flags() {
    let output = Command::new(env!("CARGO_BIN_EXE_ssh-microvm"))
        .arg("--help")
        .output()
        .expect("run ssh-microvm --help");

    assert!(output.status.success(), "--help should exit successfully");

    let stdout = String::from_utf8(output.stdout).expect("help should be valid utf-8");
    for flag in [
        "--dry-boot",
        "--listen",
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
        assert!(stdout.contains(flag), "--help should list {flag}");
    }
}
