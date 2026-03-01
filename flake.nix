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
          packages = [
            pkgs.go_1_25
            pkgs.gopls
            pkgs.gotools
            pkgs.just
            pkgs.pre-commit
            pkgs.direnv
            pkgs.nix-direnv
            pkgs.squashfsTools
            pkgs.squashfuse
          ];

          # Ensure nix-direnv's use_flake hook is available
          shellHook = ''
            source ${pkgs.nix-direnv}/share/nix-direnv/direnvrc
            pre-commit install
          '';
        };
      }
    );
}
