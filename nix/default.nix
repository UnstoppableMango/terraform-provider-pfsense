{
  a2b,
  buildGoApplication,
  gomod2nix,
  pkgs,
}:
let
  inherit (a2b.terraform) genProvider genProviderSpec scaffold;

  # Copied to the store on its own: a2b's genProvider passes the path through
  # toString, which would otherwise name a file outside the build sandbox.
  schemaFile = builtins.path {
    path = ../gen/schema.json;
    name = "schema.json";
  };

  tools = pkgs.callPackage ./tools.nix {
    inherit buildGoApplication;
  };

  openapi = pkgs.callPackage ./openapi.nix {
    inherit tools;
  };

  spec = pkgs.callPackage ./provider-spec.nix {
    inherit genProviderSpec openapi tools;
  };

  src = pkgs.callPackage ./provider-src.nix {
    inherit
      genProvider
      gomod2nix
      scaffold
      schemaFile
      tools
      ;
  };
in
buildGoApplication {
  pname = "terraform-provider-pfsense";
  version = "0.1.0";
  modules = ../gen/gomod2nix.toml;
  inherit src;

  # The Go build cache filters go.mod out of src during evaluation, and src is
  # the generated provider source, so the cache would make the whole pipeline
  # an import-from-derivation.
  disableGoCache = true;

  subPackages = [ "cmd/terraform-provider-pfsense" ];

  passthru = {
    inherit
      schemaFile
      spec
      src
      tools
      ;
  };

  ldflags = [
    "-w"
    "-s"
    "-X github.com/unstoppablemango/terraform-provider-pfsense/provider_pfsense.Version=0.1.0"
  ];
}
