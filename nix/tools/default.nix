{
  buildGoApplication,
  fetchFromGitHub,
  lib,
  runCommand,
}:
let
  fs = lib.fileset;

  # TODO: We reference this twice, once in unmango/pkgs and once here.
  # Would be nice to slurp the exact config.go the binary was built with.
  upstream = fetchFromGitHub {
    owner = "hashicorp";
    repo = "terraform-plugin-codegen-openapi";
    rev = "v0.3.0";
    sha256 = "sha256-6xI6PVlvYHwOnWjE0pKYDF/FvdomE5KydS7gBokJ2EM=";
  };

  slurp = buildGoApplication {
    pname = "slurp";
    version = "0.1.0";
    modules = ./gomod2nix.toml;

    src = fs.toSource {
      root = ./.;
      fileset = fs.unions [
        ./cmd/slurp-source
        ./go.mod
        ./go.sum
      ];
    };
  };

  configGo = runCommand "config.go" { } ''
    ${slurp}/bin/slurp-source ${upstream} $out
  '';
in
buildGoApplication {
  pname = "tools";
  version = "0.1.0";
  modules = ./gomod2nix.toml;
  passthru = { inherit configGo; };

  src = fs.toSource {
    root = ./.;
    fileset = fs.unions [
      ./cmd/gen-config
      ./cmd/patch-openapi
      ./cmd/patch-provider
      ./internal
      ./go.mod
      ./go.sum
    ];
  };
}
