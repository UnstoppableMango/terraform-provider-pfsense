{
  fetchurl,
  python3,
  runCommand,
  version ? "2.8.1",
}:
let
  raw = fetchurl {
    url = "https://github.com/pfrest/pfSense-pkg-RESTAPI/releases/download/v${version}/openapi.json";
    hash = "sha256-Va7wPj+AsrYRGBz/ZEQChJKo10oFSkxAIFNDBAHcIOI=";
  };
in
runCommand "openapi.json" { nativeBuildInputs = [ python3 ]; } ''
  python3 ${./patch-openapi.py} ${raw} $out
''
