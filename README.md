# ssh-microvm

An SSH server that boots a fresh [Firecracker](https://firecracker-microvm.github.io/)
microVM for each accepted connection and proxies shell, exec, PTY, and
window-change traffic into the guest. v1 intentionally allows one VM at a time:
additional SSH clients wait until the active VM is torn down.

> Status: **Rust v1 implementation.** Scope is intentionally small — single
> concurrent VM, cold boot only, externally managed tap networking, and no
> persistent state. See [`plan.md`](./plan.md).

## Architecture (v1)

```
ssh client ──► ssh-microvm (russh server)
                  │
                  ├── Firecracker process       (per-connection cold boot)
                  ├── Tiny UDS HTTP client      (configures Firecracker API)
                  └── russh client → guest:22   (proxies channels both ways)
```

Code layout (all under `src/`):

- `main.rs` — entrypoint; parses CLI, sets up logging, runs the server
- `config.rs` — clap-derived `Config`
- `api.rs` — hand-rolled HTTP-over-UnixStream client for the Firecracker API
- `boot.rs` — Firecracker configuration sequence and guest sshd readiness wait
- `firecracker.rs` — `VmHandle::boot` / `shutdown`, process supervision
- `lifecycle.rs` — single-VM mutex lease used by SSH sessions
- `ssh_server.rs` — russh server, auth, channel lifecycle, VM lease ownership
- `proxy.rs` — inner russh client for guest exec/shell proxying

## Requirements

- Linux with KVM (`/dev/kvm` accessible to your user)
- [Nix](https://nixos.org/) with flakes + [direnv](https://direnv.net/) (recommended), or:
  - Rust (stable, edition 2024), `just`, `squashfsTools`, `squashfuse`,
    `curl`, `mkfs.ext4`, `ssh-keygen`, `sudo`
- The `firecracker` binary on `PATH`
- `sudo` access to manage tap devices (see `just integration-sudoers`)

## Quick start

```bash
# 1. Enter dev shell (Nix flake provides Rust, just, squashfs tools, pre-commit)
direnv allow      # or: nix develop

# 2. Download the Firecracker CI kernel + Ubuntu rootfs, generate guest SSH key,
#    inject authorized_keys, build an ext4 image. Outputs into ./artifacts/.
just fetch-ubuntu

# 3. Bring up the host-side tap device (172.16.0.1/24 on tap0)
just tap-up

# 4. Run the server (listens on :2222 by default)
just run

# 5. From another shell, SSH into a freshly-booted microVM
just ssh-local
```

When you disconnect, the microVM is terminated and its state directory is
cleaned up. The next client connection cold-boots a new guest.

## Configuration

Configured via flags:

| Flag | Default | Notes |
| --- | --- | --- |
| `--dry-boot` | `false` | Boot one VM, wait for guest sshd, then tear it down |
| `microvm boot` | n/a | Subcommand equivalent of `--dry-boot` |
| `--listen` | `0.0.0.0:2222` | SSH listen address |
| `--kernel` | _(required)_ | Path to `vmlinux.bin` |
| `--rootfs` | _(required)_ | Path to ext4 rootfs |
| `--state-dir` | `/tmp/ssh-microvm` | Sockets, logs, host key, per-VM state |
| `--host-key` | `<state-dir>/ssh_host_ed25519` | Auto-generated if missing |
| `--authorized-keys` | _(none)_ | OpenSSH `authorized_keys` file |
| `--accept-any-key` | `false` | Dev only; bypasses auth |
| `--firecracker` | `firecracker` | Binary path |
| `--vcpu` / `--mem` | `1` / `512` | Guest resources |
| `--boot-args` | `console=ttyS0 reboot=k panic=1 pci=off` | Kernel cmdline |
| `--guest-user` | `root` | SSH user inside the guest |
| `--guest-key` | _(required)_ | Private key used to reach guest sshd |
| `--guest-ip` / `--host-ip` | `172.16.0.2` / `172.16.0.1` | Tap network |
| `--tap-name` | `tap0` | Host tap device (created externally) |
| `--boot-timeout` | `30s` | How long to wait for guest sshd |
| `--grace-stop` | `2s` | Time before SIGKILL on shutdown |

## Common tasks (Justfile)

```bash
just                    # list recipes
just build              # cargo build --release
just run                # run with artifacts/vmlinux.bin + artifacts/ubuntu.ext4
just fmt / fmt-check    # cargo fmt
just lint               # cargo clippy -D warnings
just check              # cargo check --all-targets
just test               # cargo test
just fetch-ubuntu       # download + assemble Ubuntu rootfs
just tap-up / tap-down  # manage tap0 / 172.16.0.1
ssh-microvm --dry-boot ... # boot one VM and verify guest sshd readiness
ssh-microvm ... microvm boot # subcommand form of --dry-boot
just ssh-local          # ssh into the running server
```

## Development

- Pre-commit hooks (`cargo fmt --check`, `cargo check`) installed by the dev shell.
- Rust toolchain comes from `nixpkgs` unstable (no overlay pin yet).
- Full local validation: `cargo fmt --all -- --check`, `cargo test --all`, and
  `cargo clippy --all-targets -- -D warnings`.
- The smoke integration test is gated by `SSH_MICROVM_KERNEL`,
  `SSH_MICROVM_ROOTFS`, and `SSH_MICROVM_GUEST_KEY`. When all three are set,
  `cargo test --test smoke` boots the binary, runs `echo hello` through SSH, and
  asserts the proxied command succeeds.

## License

Licensed under the Apache License, Version 2.0. See [`LICENSE`](./LICENSE).
