{
  config,
  buildProviderSpec,
  openapi,
}:
buildProviderSpec {
  name = "terraform-provider-pfsense";
  src = openapi;
  inherit config;
}
