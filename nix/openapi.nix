{
  fetchurl,
  python3,
  runCommand,
  tools,
  version ? "2.8.1",
}:
let
  raw = fetchurl {
    url = "https://github.com/pfrest/pfSense-pkg-RESTAPI/releases/download/v${version}/openapi.json";
    hash = "sha256-Va7wPj+AsrYRGBz/ZEQChJKo10oFSkxAIFNDBAHcIOI=";
  };
in
runCommand "openapi.json" { } ''
  ${tools}/bin/patch-openapi ${raw} $out
''
