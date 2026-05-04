# Justfile for ssh-microvm (Rust)

set shell := ["bash", "-lc"]

KERNEL := "artifacts/vmlinux.bin"
ROOTFS := "artifacts/ubuntu.ext4"
GUEST_KEY := "artifacts/ubuntu.id_rsa"
TAP_NAME := "tap0"
HOST_IP := "172.16.0.1"

default:
    @just --list

# --- Rust ---

build:
    cargo build --release

run:
    cargo run -- \
        --kernel {{KERNEL}} \
        --rootfs {{ROOTFS}} \
        --guest-key {{GUEST_KEY}} \
        --accept-any-key

check:
    cargo check --all-targets

fmt:
    cargo fmt --all

fmt-check:
    cargo fmt --all -- --check

lint:
    cargo clippy --all-targets -- -D warnings

test:
    cargo test --all

setup:
    pre-commit install

# --- Artifacts (unchanged from the Go version) ---

fetch-ubuntu:
    #!/usr/bin/env bash
    set -euxo pipefail
    for cmd in curl squashfuse fusermount ssh-keygen mkfs.ext4 sudo; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            echo "missing dependency: $cmd" >&2
            exit 1
        fi
    done
    mkdir -p artifacts
    ARCH="$(uname -m)"
    release_url="https://github.com/firecracker-microvm/firecracker/releases"
    latest_version="$(basename "$(curl -fsSLI -o /dev/null -w '%{url_effective}' "${release_url}/latest")")"
    CI_VERSION="${latest_version%.*}"
    latest_kernel_key="$(curl -fsSL "http://spec.ccfc.min.s3.amazonaws.com/?prefix=firecracker-ci/${CI_VERSION}/${ARCH}/vmlinux-&list-type=2" | grep -oE "firecracker-ci/${CI_VERSION}/${ARCH}/vmlinux-[0-9]+\\.[0-9]+\\.[0-9]{1,3}" | sort -V | tail -1)"
    curl -fsSL -o artifacts/vmlinux.bin "https://s3.amazonaws.com/spec.ccfc.min/${latest_kernel_key}"
    latest_ubuntu_key="$(curl -fsSL "http://spec.ccfc.min.s3.amazonaws.com/?prefix=firecracker-ci/${CI_VERSION}/${ARCH}/ubuntu-&list-type=2" | grep -oE "firecracker-ci/${CI_VERSION}/${ARCH}/ubuntu-[0-9]+\\.[0-9]+\\.squashfs" | sort -V | tail -1)"
    ubuntu_base="$(basename "${latest_ubuntu_key}")"
    ubuntu_version="${ubuntu_base#ubuntu-}"
    ubuntu_version="${ubuntu_version%.squashfs}"
    curl -fsSL -o "artifacts/ubuntu-${ubuntu_version}.squashfs.upstream" "https://s3.amazonaws.com/spec.ccfc.min/${latest_ubuntu_key}"
    squash_mount="$(mktemp -d /tmp/ssh-microvm-squashfs-mount-XXXXXXXX)"
    squash_dir="$(mktemp -d /tmp/ssh-microvm-squashfs-XXXXXXXX)"
    squashfuse "artifacts/ubuntu-${ubuntu_version}.squashfs.upstream" "$squash_mount"
    cp -a "$squash_mount/." "$squash_dir/"
    fusermount -u "$squash_mount"
    rmdir "$squash_mount"
    key_path="artifacts/ubuntu-${ubuntu_version}.id_rsa"
    if [ ! -f "$key_path" ]; then
        ssh-keygen -f "$key_path" -N ""
    fi
    ln -sf "ubuntu-${ubuntu_version}.id_rsa" artifacts/ubuntu.id_rsa
    ln -sf "ubuntu-${ubuntu_version}.id_rsa.pub" artifacts/ubuntu.id_rsa.pub
    mkdir -p "$squash_dir/root/.ssh"
    cp -v artifacts/ubuntu-${ubuntu_version}.id_rsa.pub "$squash_dir/root/.ssh/authorized_keys"
    sudo chown -R root:root "$squash_dir"
    truncate -s 1G "artifacts/ubuntu-${ubuntu_version}.ext4"
    sudo mkfs.ext4 -d "$squash_dir" -F "artifacts/ubuntu-${ubuntu_version}.ext4"
    sudo rm -rf "$squash_dir"
    ln -sf "ubuntu-${ubuntu_version}.ext4" artifacts/ubuntu.ext4

# --- Networking (host tap) ---

tap-up:
    #!/usr/bin/env bash
    set -euo pipefail
    tap="{{TAP_NAME}}"
    host_ip="{{HOST_IP}}"
    if ip link show "$tap" >/dev/null 2>&1; then
        echo "tap exists: $tap"
    else
        sudo ip tuntap add dev "$tap" mode tap
    fi
    if ! ip addr show dev "$tap" | grep -q "${host_ip}/24"; then
        sudo ip addr add "${host_ip}/24" dev "$tap"
    fi
    sudo ip link set "$tap" up

tap-down:
    #!/usr/bin/env bash
    set -euo pipefail
    tap="{{TAP_NAME}}"
    if ip link show "$tap" >/dev/null 2>&1; then
        sudo ip link del "$tap"
    else
        echo "tap missing: $tap"
    fi

integration-sudoers:
    #!/usr/bin/env bash
    set -euo pipefail
    user="${SUDO_USER:-${USER:-$LOGNAME}}"
    cat <<'EOF'
    Add the following to your sudoers (via `sudo visudo`), updating USER if needed:
    EOF
    cat <<EOF
    ${user} ALL=(root) NOPASSWD: /usr/sbin/ip, /sbin/ip, /usr/bin/ip
    EOF

# --- Connect ---

ssh-local:
    ssh -t -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \
        -o IdentitiesOnly=yes -i artifacts/ubuntu.id_rsa \
        localhost -p 2222
