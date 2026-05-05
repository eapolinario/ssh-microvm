# ssh-microvm plan (Rust rewrite)

## v1 scope (one sentence)

An SSH server that, for each accepted connection, cold-boots a Firecracker
microVM, waits for guest sshd, and bidirectionally proxies the SSH session
(shell / exec / PTY).

## Constraints

- Single concurrent VM. Subsequent SSH connections wait on a tokio mutex.
- Host-side tap (`tap0`, 172.16.0.1/24) is created externally via `just tap-up`.
- No persistent state. No SQLite. No user store. No audit log.
- Auth: `--authorized-keys <path>` (OpenSSH format) or `--accept-any-key` (dev).
- Cold-boot only. No snapshots.

## Step-by-step execution

- [x] **1. Cleanup + scaffolding** — delete Go tree, write `Cargo.toml`,
      flake.nix (Rust toolchain from nixpkgs unstable), Justfile, README,
      LICENSE (Apache-2.0), `src/main.rs` + `src/config.rs` stubs.
- [x] **2. Config + CLI** — full `Config` struct via clap derive; `--help`
      output documents every flag; `cargo run -- --help` works.
- [x] **3. Firecracker UDS HTTP client** (`src/api.rs`) — hand-rolled PUT/GET
      over `tokio::net::UnixStream`. Plus `--dry-boot` / `microvm boot`
      subcommand that spawns Firecracker, configures it, starts the VM,
      waits for guest sshd:22, and exits cleanly on Ctrl-C.
- [x] **4. VmHandle** (`src/firecracker.rs`) — `boot()` + `shutdown()` with
      SendCtrlAltDel → grace → SIGKILL → state-dir cleanup.
- [x] **5. Outer SSH server** (`src/ssh_server.rs`) — russh server, host-key
      load (auto-generate if missing), public-key auth against
      `--authorized-keys` or `--accept-any-key`. No VM yet; reply "ok" to
      exec to validate the SSH stack.
- [x] **6. VM lifecycle on connection** — global `tokio::sync::Mutex`; boot
      VM lazily on first channel open; shut down on disconnect.
- [x] **7. Inner SSH client + proxy** (`src/proxy.rs`) — russh client dials
      `guest_ip:22` with `--guest-key`; proxy `exec` end-to-end.
- [x] **8. PTY / shell / window-change** — full interactive proxy.
- [x] **9. Smoke integration test** — gated by `SSH_MICROVM_KERNEL`,
      `SSH_MICROVM_ROOTFS`, `SSH_MICROVM_GUEST_KEY`. Boots the binary, runs
      `echo hello`, asserts exit 0 and stdout.
- [x] **10. README + plan polish.**

## Backlog (post-v1)

- [ ] Multiple concurrent VMs with per-slot tap + /30 allocation
- [ ] In-binary tap creation (drop the `just tap-up` requirement)
- [ ] Snapshot-backed warm pool, target <200ms accept→shell
- [ ] Per-session writable overlay on top of a shared base disk
- [ ] Debian rootfs build path
- [ ] Authorized-keys hot reload
- [ ] vsock console and structured boot log capture
