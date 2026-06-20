{
  genProvider,
  go,
  gomod2nix,
  input,
  lib,
  runCommand,
  scaffold,
  stdenv,
  symlinkJoin,
  writeShellApplication,
}:
let
  goPackage = "github.com/unstoppablemango/terraform-provider-pfsense";
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

  goSrc = symlinkJoin {
    name = "go-src";
    paths = [
      (genProvider {
        name = "terraform-provider-pfsense";
        inherit input;
      })
      (runCommand "go.mod" { } ''
        mkdir -p $out && cd $out
        ${go}/bin/go mod init ${goPackage}
      '')
      (scaffold {
        command = "provider";
        name = "pfSense";
        scaffoldName = "pfsense";

        # a2b scaffold does not pre-create $out; preRun hook does it
        env.preRun = "mkdir -p $out";
      })
    ]
    ++ map toScaffold schema.resources;
  };

  goModPatch = writeShellApplication {
    name = "go.mod.patch.sh";
    runtimeInputs = [ go ];

    text = ''
      go -C ${goSrc} mod tidy -diff > "$1"
    '';
  };

  patched = stdenv.mkDerivation {
    name = "patched-src";
    src = null;
    dontUnpack = true;

    prePatch = ''
      cp ${goSrc}/go.mod go.mod
    '';

    patches = [ ./go.mod.patch ];

    buildPhase = ''
      mkdir -p $out
      cp go.mod $out/go.mod
      cp go.sum $out/go.sum
    '';
  };

  gomod2nixToml = writeShellApplication {
    name = "gomod2nix.toml.sh";
    runtimeInputs = [ gomod2nix ];

    text = ''
      set -x
      gomod2nix generate --dir "${patched}" --outdir "$1"
    '';
  };
in
patched.overrideAttrs (oldAttrs: {
  passthru = {
    inherit goModPatch gomod2nixToml;
  };
})
