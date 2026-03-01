# ssh-microvm plan

## Done

- [x] Define SQLite schema + data access layer — minimal tables for users/keys/sessions/vms/audit; Go migrations + repo layer
- [x] Firecracker launch/teardown stub — process wrapper with config rendering, socket management, termination on disconnect
- [x] Implement embedded SSH server skeleton — Go SSH server, host keys, auth hooks, session lifecycle wiring
- [x] App config + minimal CLI — config struct, env/flags parsing, app entrypoint wiring
- [x] Bootstrap Go control plane skeleton — wires SSH server → Firecracker manager → SQLite store
- [x] Fetch Ubuntu kernel/rootfs artifacts + set Justfile defaults — download Firecracker quickstart artifacts, Justfile defaults
- [x] Add integration test for SSH→VM startup — starts ssh-microvm, performs SSH handshake, asserts VM startup/teardown

## In progress

- [ ] Implement SSH proxy to guest via VM networking — tap networking + guest IP, ensure sshd in guest, proxy outer SSH session to guest:22
- [ ] Run ssh-microvm binary locally — build and run with Ubuntu artifacts, verify Firecracker starts and SSH proxy works end-to-end

## Backlog

- [ ] Wire SSH session to guest — SSH I/O wiring to the microVM (vsock/serial console or forwarded port) so interactive shells run inside the guest
- [ ] Add Debian rootfs support — document and automate building a Debian rootfs ext4 image; update config to support rootfs selection
