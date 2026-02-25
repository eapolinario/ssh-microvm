# Justfile for ssh-microvm

set shell := ["bash", "-lc"]

KERNEL := "artifacts/vmlinux.bin"
ROOTFS := "artifacts/ubuntu.ext4"

default:
    @just --list

build:
    CGO_ENABLED=0 go build ./cmd/ssh-microvm

run:
    CGO_ENABLED=0 go run ./cmd/ssh-microvm --kernel {{KERNEL}} --rootfs {{ROOTFS}}

run-args KERNEL ROOTFS:
    CGO_ENABLED=0 go run ./cmd/ssh-microvm --kernel {{KERNEL}} --rootfs {{ROOTFS}}

fetch-ubuntu:
    #!/usr/bin/env bash
    set -euxo pipefail
    for cmd in curl unsquashfs ssh-keygen mkfs.ext4 sudo; do
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
    squash_dir="$(mktemp -d /tmp/ssh-microvm-squashfs-XXXXXXXX)"
    unsquashfs -d "$squash_dir" "artifacts/ubuntu-${ubuntu_version}.squashfs.upstream"
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

setup-hooks:
    git config core.hooksPath .githooks

fmt:
    gofmt -w cmd internal

lint:
    go vet ./...

tidy:
    go mod tidy

test:
    CGO_ENABLED=0 go test ./...
