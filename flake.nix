{
  description = "ssh-microvm development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    nix-direnv.url = "github:nix-community/nix-direnv";
  };

  outputs = { self, nixpkgs, flake-utils, nix-direnv }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        devShells.default = pkgs.mkShell {
          # Rust toolchain comes from nixpkgs (no rust-overlay / fenix pin yet).
          # If we ever need a pinned toolchain, swap to rust-overlay here.
          packages = [
            pkgs.rustc
            pkgs.cargo
            pkgs.rustfmt
            pkgs.clippy
            pkgs.rust-analyzer

            pkgs.just
            pkgs.pre-commit
            pkgs.direnv
            pkgs.nix-direnv

            # Used by `just fetch-ubuntu` to unpack the Firecracker CI rootfs.
            pkgs.squashfsTools
            pkgs.squashfuse
          ];

          shellHook = ''
            source ${pkgs.nix-direnv}/share/nix-direnv/direnvrc
            pre-commit install >/dev/null 2>&1 || true
          '';
        };
      }
    );
}
