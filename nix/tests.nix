{
  buildGoApplication,
  lib,
  pkgs,
  providerBin,
}:
let
  fs = lib.fileset;

  # opentofu is Apache-licensed; wrap as "terraform" since terraform-plugin-testing
  # looks for that binary name in PATH.
  terraform = pkgs.writeShellScriptBin "terraform" ''
    exec ${pkgs.opentofu}/bin/tofu "$@"
  '';

  testBin = buildGoApplication {
    pname = "pfsense-integration-tests";
    version = "0.1.0";
    modules = ../test/gomod2nix.toml;

    src = fs.toSource {
      root = ../test;
      fileset = fs.unions [
        ../test/go.mod
        ../test/go.sum
        ../test/integration
        ../test/mock
      ];
    };

    doCheck = false;

    buildPhase = ''
      go test -c -o integration-tests ./integration/
    '';

    installPhase = ''
      mkdir -p $out/bin
      cp integration-tests $out/bin/
    '';
  };
in
pkgs.writeShellApplication {
  name = "test";
  runtimeInputs = [ terraform testBin ];
  text = ''
    export PFSENSE_PROVIDER_BINARY="${providerBin}/bin/terraform-provider-pfsense"
    TF_ACC=1 integration-tests "$@"
  '';
}
