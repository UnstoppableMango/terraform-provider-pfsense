{
  pkgs,
  providerBin,
}:
let
  terraform = pkgs.writeShellScriptBin "terraform" ''
    exec ${pkgs.opentofu}/bin/tofu "$@"
  '';
in
pkgs.writeShellApplication {
  name = "test";
  runtimeInputs = [
    terraform
    pkgs.go
    pkgs.git
  ];
  text = ''
    export PFSENSE_PROVIDER_BINARY="${providerBin}/bin/terraform-provider-pfsense"
    cd "$(git rev-parse --show-toplevel)/test"
    TF_ACC=1 go test "$@" ./...
  '';
}
