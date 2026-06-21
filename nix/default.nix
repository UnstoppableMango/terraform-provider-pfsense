{
  a2b,
  buildGoApplication,
  gomod2nix,
  pkgs,
}:
let
  inherit (a2b.terraform) genProvider genProviderSpec scaffold;

  tools = pkgs.callPackage ./tools {
    inherit buildGoApplication;
  };

  openapi = pkgs.callPackage ./openapi.nix {
    inherit tools;
  };

  spec = pkgs.callPackage ./provider-spec.nix {
    inherit genProviderSpec openapi tools;
  };

  src = pkgs.callPackage ./provider-src.nix {
    inherit genProvider gomod2nix scaffold tools;
    schemaFile = spec;
  };
in
buildGoApplication {
  pname = "terraform-provider-pfsense";
  version = "0.1.0";
  modules = ./gomod2nix.toml;
  inherit src;

  subPackages = [ "cmd/terraform-provider-pfsense" ];

  passthru = { inherit spec src tools; };

  ldflags = [
    "-w"
    "-s"
  ];
}
