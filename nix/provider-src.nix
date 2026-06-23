{
  applyPatches,
  genProvider,
  go,
  gomod2nix,
  schemaFile,
  lib,
  runCommand,
  scaffold,
  stdenv,
  symlinkJoin,
  tools,
  writeShellApplication,
}:
let
  fs = lib.fileset;
  goPackage = "github.com/unstoppablemango/terraform-provider-pfsense";
  schema = builtins.fromJSON (builtins.readFile schemaFile);

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

      env.postRun = ''
        mkdir -p $out/${package}
        mv $out/${name}.go $out/${package}/
      '';
    };

  patchedProvider = runCommand "provider_pfsense" { } ''
    mkdir -p $out/provider_pfsense
    ${tools}/bin/patch-provider \
      ${../provider_pfsense/provider.go} \
      ${schemaFile} \
      > $out/provider_pfsense/provider.go
  '';

  goSrc = symlinkJoin {
    name = "go-src";
    paths = [
      (fs.toSource {
        root = ../.;
        fileset = ../cmd;
      })
      patchedProvider
      (genProvider {
        name = "terraform-provider-pfsense";
        input = schemaFile;
      })
      (runCommand "go.mod" { } ''
        mkdir -p $out && cd $out
        ${go}/bin/go mod init ${goPackage}
      '')
      # This scaffolds the same thing we're generating with patchedProvider
      # TODO: Retarget patches to the scaffolding output
      # (scaffold {
      #   command = "provider";
      #   name = "pfSense";
      #   scaffoldName = "pfsense";

      #   # a2b scaffold does not pre-create $out; preRun hook does it
      #   env.preRun = "mkdir -p $out";
      # })
    ]
    ++ map toScaffold schema.resources;
  };

  goModPatch = writeShellApplication {
    name = "go.mod.patch.sh";
    runtimeInputs = [ go ];

    text = ''
      go -C ${goSrc} mod tidy -diff >"$1" || true
    '';
  };

  patched = applyPatches {
    src = runCommand "deref-symlinks" { } ''
      cp -rL ${goSrc} $out
    '';
    patches = [ ./go.mod.patch ];
  };

  gomod2nixToml = writeShellApplication {
    name = "gomod2nix.toml.sh";
    runtimeInputs = [ gomod2nix ];

    text = ''
      gomod2nix generate --dir "${patched}" --outdir "$1"
    '';
  };
in
patched.overrideAttrs (oldAttrs: {
  passthru = {
    inherit goModPatch gomod2nixToml;
  };
})
