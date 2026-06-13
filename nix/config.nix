{
  openapi,
  runCommand,
  tools,
}:
runCommand "config.yaml" { } ''
  ${tools}/bin/gen-config ${openapi} $out
''
