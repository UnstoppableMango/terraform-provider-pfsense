{
  genOpenapi,
  openapi,
}:
genOpenapi {
  name = "terraform-provider-pfsense";
  src = openapi;
  config = ../generator_config.yml;
}
