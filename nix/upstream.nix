{
  fetchFromGitHub,
  runCommand,
  tools,
}:
let
  src = fetchFromGitHub {
    owner = "hashicorp";
    repo = "terraform-plugin-codegen-openapi";
    rev = "v0.3.0";
    sha256 = "sha256-6xI6PVlvYHwOnWjE0pKYDF/FvdomE5KydS7gBokJ2EM=";
  };
in
runCommand "config.go" { } ''
  ${tools}/bin/slurp-source ${src} $out
''
