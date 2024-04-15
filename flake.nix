{
  description = "Banoffee flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      devShell = pkgs.mkShell { 
        NIX_CONFIG = "experimental-features = nix-command flakes";
        nativeBuildInputs = with pkgs; [ nix git nodejs go golangci-lint gopls ];
        buildInputs = with pkgs.python3Packages; [
          python
          requests
        ];
      };
  });
}
