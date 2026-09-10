# terraform-provider-pfsense

[![CI](https://github.com/UnstoppableMango/terraform-provider-pfsense/actions/workflows/ci.yml/badge.svg)](https://github.com/UnstoppableMango/terraform-provider-pfsense/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/unstoppablemango/terraform-provider-pfsense.svg)](https://pkg.go.dev/github.com/unstoppablemango/terraform-provider-pfsense)
[![Go Version](https://img.shields.io/github/go-mod/go-version/UnstoppableMango/terraform-provider-pfsense)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Nix Flake](https://img.shields.io/badge/nix-flake-5277C3?logo=nixos)](flake.nix)
[![GitHub last commit](https://img.shields.io/github/last-commit/UnstoppableMango/terraform-provider-pfsense)](https://github.com/UnstoppableMango/terraform-provider-pfsense/commits/main)

A [Terraform](https://www.terraform.io/) provider for [pfSense](https://www.pfsense.org/), generated entirely from the [pfSense REST API](https://github.com/pfrest/pfSense-pkg-RESTAPI) OpenAPI specification via a [Nix](https://nixos.org/)-driven pipeline.
No provider code is written by hand.

## How it works

The provider is produced by a multi-stage code generation pipeline that transforms the upstream pfSense OpenAPI spec into provider Go source code.

```mermaid
flowchart TD
    A["pfSense REST API\n(OpenAPI JSON)"]
    B["patch-openapi\n(flatten allOf)"]
    C["Patched OpenAPI spec\n(openapi.json)"]
    D["gen-config\n(build generator config)"]
    E["config.yaml\n(tfplugingen-openapi config)"]
    F["tfplugingen-openapi\nvia a2b genProviderSpec"]
    G["schema.json\n(provider schema)"]
    H["a2b genProvider + scaffold\n+ gen-main + patch-provider"]
    I["Provider Go source\n(generated)"]
    J["slurp-source\n(extract config.go AST)"]
    K["hashicorp/terraform-plugin-codegen-openapi\n(upstream source)"]
    L["internal/config/config.go\n(Config struct)"]
    M["Go compiler"]
    N["terraform-provider-pfsense\n(final binary)"]

    K -->|"nix/tools.nix"| J
    J --> L
    L --> D
    A -->|"nix/openapi.nix"| B
    B --> C
    C -->|"nix/provider-spec.nix"| D
    D --> E
    E --> F
    C --> F
    F --> G
    G -->|"nix/provider-src.nix"| H
    H --> I
    I -->|"nix/default.nix"| M
    M --> N
```

### Pipeline stages

| Nix derivation | Tool | Input | Output |
|---|---|---|---|
| `nix/tools.nix` | Go compiler + `slurp-source` | `cmd/` Go source + [`hashicorp/terraform-plugin-codegen-openapi`](https://github.com/hashicorp/terraform-plugin-codegen-openapi) | CLI tools binary + `internal/config/config.go` |
| `nix/openapi.nix` | `patch-openapi` | pfSense REST API [OpenAPI release](https://github.com/pfrest/pfSense-pkg-RESTAPI/releases) | Patched `openapi.json` |
| `nix/provider-spec.nix` | `gen-config` + [`tfplugingen-openapi`](https://github.com/hashicorp/terraform-plugin-codegen-openapi) via [`a2b`](https://github.com/UnstoppableMango/a2b) | Patched OpenAPI spec | `config.yaml` → `schema.json` |
| `nix/provider-src.nix` | [`tfplugingen-sdk`](https://github.com/hashicorp/terraform-plugin-codegen-sdk) via [`a2b`](https://github.com/UnstoppableMango/a2b) + `gen-main` + `patch-provider` | `schema.json` | Generated provider Go source |
| `nix/default.nix` | Go compiler | Provider Go source | `terraform-provider-pfsense` binary |

## Requirements

- [Nix](https://nixos.org/) with flakes enabled
- [direnv](https://direnv.net/) (optional, for automatic dev shell activation)

## Usage

```bash
# Build the default output (provider binary)
nix build

# Build individual pipeline stages
nix build .#tools   # CLI tools binary
nix build .#src     # generated provider Go source
nix build .#bin     # provider binary (same as default)
```

## Development

```bash
# Enter dev shell (or use direnv)
nix develop

# Run tests
go test ./...

# Lint / format
nix flake check

# Update Go dependencies
go mod tidy
gomod2nix generate --outdir ./nix

# Update Nix flake inputs
make update
```

### Adding a resource

Edit `cmd/gen-config/config.go` -> `ConfigFor()` to add an entry to the `Resources` map with the appropriate OpenAPI path and method for each CRUD operation.
The downstream Nix pipeline picks up the change on the next `nix build`.

## Key dependencies

- [pb33f/libopenapi](https://github.com/pb33f/libopenapi) - OpenAPI v3 parsing and bundling
- [hashicorp/terraform-plugin-codegen-openapi](https://github.com/hashicorp/terraform-plugin-codegen-openapi) - generates provider schema from OpenAPI spec
- [hashicorp/terraform-plugin-codegen-sdk](https://github.com/hashicorp/terraform-plugin-codegen-sdk) - generates provider Go code from schema
- [UnstoppableMango/a2b](https://github.com/UnstoppableMango/a2b) - Nix functions wrapping the HashiCorp codegen toolchain
- [nix-community/gomod2nix](https://github.com/nix-community/gomod2nix) - reproducible Go builds in Nix
- [pfrest/pfSense-pkg-RESTAPI](https://github.com/pfrest/pfSense-pkg-RESTAPI) - upstream pfSense REST API and OpenAPI spec

## License

[MIT](LICENSE)
