{
  config,
  genProviderSpec,
  openapi,
}:
genProviderSpec {
  name = "schema.json";
  inherit config;
  src = openapi;
}
