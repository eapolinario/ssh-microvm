# ssh-microvm

Minimal control plane that starts a Firecracker microVM when an SSH client connects.

## Quick start (Ubuntu artifacts)

```bash
just fetch-ubuntu
just run
```

This downloads the latest Firecracker CI Ubuntu kernel + rootfs, builds an ext4 rootfs, and starts the SSH server.

### Requirements for `fetch-ubuntu`

- `curl`
- `unsquashfs`
- `ssh-keygen`
- `mkfs.ext4` (requires sudo)

## Run with explicit paths

```bash
just run-args /path/to/vmlinux.bin /path/to/rootfs.ext4
```

## Defaults

- Listen: `:2222`
- State dir: `/tmp/ssh-microvm-<random-id>` (generated on startup)
- DB: `<state-dir>/db.sqlite`
- Auth mode: `auto-enroll` (accepts any public key and enrolls it on first connect)
- Kernel: `artifacts/vmlinux.bin`
- Rootfs: `artifacts/ubuntu.ext4`
- Guest SSH key: `artifacts/ubuntu.id_rsa` (created by `just fetch-ubuntu`)
- Guest IP: `172.16.0.2` (host tap: `172.16.0.1`)

## Behavior on disconnect

When the SSH session ends, the microVM is terminated (SIGTERM, then hard kill after a short grace period).

## SSH forwarding to guest

This uses a host-side SSH client to connect to the guest’s `sshd` and bridges the SSH session I/O.
It requires:
- `sudo` access to create a tap device (`sudo ip tuntap ...`).
- The guest rootfs to include `sshd` with root key-based login enabled.

## Auth modes

- `auto-enroll` (default): accept any key, store it in SQLite after successful handshake.
- `known-keys`: only allow keys already present in SQLite.

## Notes

- Networking is minimal and assumes a single VM at a time (static guest IP).
- This v1 skeleton does not yet handle multi-VM IP allocation or NAT.
