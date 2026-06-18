{
  genProvider,
  input,
  lib,
  scaffold,
  symlinkJoin,
}:
let
  schema = builtins.fromJSON (builtins.readFile input);

  toScaffold =
    resource:
    let
      name = "${resource.name}_resource";
      package = "resource_${name}";
    in
    scaffold {
      command = "resource";
      inherit name package;
      scaffoldName = lib.strings.toSentenceCase resource.name;

      # a2b scaffold does not pre-create $out; preRun hook does it
      env.preRun = "mkdir -p $out";

      env.postRun = ''
        mkdir -p $out/${package}
        mv $out/${name}.go $out/${package}/
      '';
    };
in
symlinkJoin {
  name = "provider-src";
  paths = [
    (genProvider {
      name = "terraform-provider-pfsense";
      inherit input;
    })
    (scaffold {
      command = "provider";
      name = "pfSense";
      scaffoldName = "pfsense";

      # a2b scaffold does not pre-create $out; preRun hook does it
      env.preRun = "mkdir -p $out";
    })
  ]
  ++ map toScaffold schema.resources;
}
