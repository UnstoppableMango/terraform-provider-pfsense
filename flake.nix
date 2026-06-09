{
  description = "A Nix flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    systems.url = "github:nix-systems/default";

    globset = {
      url = "github:pdtpartners/globset";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };

    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };

    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.inputs.systems.follows = "systems";
    };

    mangonix = {
      url = "github:UnstoppableMango/nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.systems.follows = "systems";
      inputs.gomod2nix.follows = "gomod2nix";
      inputs.flake-parts.follows = "flake-parts";
      inputs.treefmt-nix.follows = "treefmt-nix";
    };
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import inputs.systems;
      imports = with inputs; [ treefmt-nix.flakeModule ];

      perSystem =
        { pkgs, inputs', ... }:
        let
          inherit (inputs'.mangonix.legacyPackages) terraformTools;
          inherit (inputs'.gomod2nix.legacyPackages) buildGoApplication gomod2nix;
          inherit (inputs.globset.lib) globs;

          tools = pkgs.callPackage ./nix/tools.nix { inherit buildGoApplication globs; };
          openapi = pkgs.callPackage ./nix/openapi.nix { inherit tools; };

          generator = pkgs.callPackage ./nix {
            inherit (terraformTools) genOpenapi;
            inherit openapi;
          };
        in
        {
          packages = {
            inherit tools openapi;
            default = generator;
          };

          devShells.default = pkgs.mkShellNoCC {
            packages = with pkgs; [
              direnv
              go
              gomod2nix
              gopls
              ginkgo
              gnumake
              nixfmt
            ];

            GO = "${pkgs.go}/bin/go";
            GOMOD2NIX = "${gomod2nix}/bin/gomod2nix";
          };

          treefmt.programs = {
            actionlint.enable = true;
            nixfmt.enable = true;
            gofmt.enable = true;
          };
        };
    };
}
