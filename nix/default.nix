{
  config,
  genOpenapi,
  openapi,
}:
genOpenapi {
  name = "terraform-provider-pfsense";
  src = openapi;
  inherit config;
}
