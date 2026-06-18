{
  a2b,
  buildGoApplication,
  globs,
  pkgs,
  tools,
}:
let
  inherit (a2b.terraform) genProvider genProviderSpec scaffold;

  openapi = pkgs.callPackage ./openapi.nix { inherit tools; };

  spec = pkgs.callPackage ./provider-spec.nix {
    inherit genProviderSpec openapi tools;
  };

  src = pkgs.callPackage ./provider-src.nix {
    inherit genProvider scaffold;
    input = spec;
  };
in
buildGoApplication {
  pname = "terraform-provider-pfsense";
  version = "0.1.0";
  modules = ./gomod2nix.toml;
  inherit src;

  subPackages = [ "cmd/terraform-provider-pfsense" ];

  passthru = { inherit spec src; };

  ldflags = [
    "-w"
    "-s"
  ];
}
