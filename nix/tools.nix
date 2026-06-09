{
  buildGoApplication,
  lib,
  globs,
}:
let
  fs = lib.fileset;
in
buildGoApplication {
  pname = "tfgen-pfsense-tools";
  version = "0.1.0";
  modules = ./gomod2nix.toml;

  src = fs.toSource {
    root = ../.;
    fileset = globs ../. [
      "go.mod"
      "go.sum"
      "**/*.go"
    ];
  };
}
