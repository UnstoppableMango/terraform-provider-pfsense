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

    a2b = {
      url = "github:UnstoppableMango/a2b?ref=even-more-tf";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.systems.follows = "systems";
      inputs.flake-parts.follows = "flake-parts";
      inputs.gomod2nix.follows = "gomod2nix";
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
          inherit (inputs'.gomod2nix.legacyPackages) buildGoApplication gomod2nix;
          inherit (inputs.globset.lib) globs;

          a2b = inputs'.a2b.legacyPackages.lib;

          tools = pkgs.callPackage ./nix/tools.nix { inherit buildGoApplication globs; };
          upstream = pkgs.callPackage ./nix/upstream.nix { inherit tools; };
          openapi = pkgs.callPackage ./nix/openapi.nix { inherit tools; };
          config = pkgs.callPackage ./nix/config.nix { inherit tools openapi; };

          spec = pkgs.callPackage ./nix/provider-spec.nix {
            inherit (a2b.terraform) genProviderSpec;
            inherit config openapi;
          };

          provider = pkgs.callPackage ./nix {
            inherit (a2b.terraform) genProvider;
            input = spec;
          };
        in
        {
          packages = {
            inherit
              tools
              openapi
              upstream
              config
              spec
              provider
              ;

            default = provider;
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
