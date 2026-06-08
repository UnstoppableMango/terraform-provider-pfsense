{
  buildGoApplication,
  lib,
}:
buildGoApplication {
  pname = "tfgen-pfsense-tools";
  version = "0.1.0";

  src = lib.cleanSource ./..;
  modules = ./gomod2nix.toml;
}
