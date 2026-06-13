{
  config,
  genProviderSpec,
  openapi,
}:
genProviderSpec {
  name = "schema.json";
  inherit config;
  openapi-spec = openapi;
}
