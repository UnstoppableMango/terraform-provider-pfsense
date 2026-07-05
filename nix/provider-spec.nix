{
  config,
  genProviderSpec,
  openapi,
  runCommand,
  tools,
}:
genProviderSpec {
  name = "schema.json";
  openapi-spec = openapi;

  config = runCommand "config.yaml" { } ''
    ${tools}/bin/gen-plugingen-config ${openapi} $out
  '';
}
