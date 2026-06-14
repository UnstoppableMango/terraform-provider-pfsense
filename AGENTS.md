# AGENTS.md

This file provides guidance to AI agents when working with code in this repository.

## Commands

```bash
# Build (uses nix)
make build          # or: nix build

# Run tests
go test ./...
go test ./pkg/...   # single package

# Lint / format
nix flake check     # runs actionlint, nixfmt, gofmt via treefmt

# Update go dependencies
go mod tidy                              # updates go.sum
gomod2nix generate --outdir ./nix       # sync nix/gomod2nix.toml

# Update nix flake inputs
make update         # or: nix flake update

# Build individual nix packages
nix build .#tools       # CLI tools binary
nix build .#openapi     # patched pfSense OpenAPI spec
nix build .#config      # generated tfplugingen-openapi config YAML
nix build .#spec        # generated provider schema JSON
nix build .#provider    # generated provider Go source (default)
```

## Architecture

This project generates a Terraform provider for pfSense entirely from the pfSense REST API's OpenAPI spec, using a Nix-driven pipeline. No provider code is written by hand — it is all generated.

### Code generation pipeline (nix build order)

1. **`nix/tools.nix`** — builds the CLI tools in this repo (`cmd/`) into a single binary via `gomod2nix`.
2. **`nix/upstream.nix`** — fetches `internal/config/parse.go` from `hashicorp/terraform-plugin-codegen-openapi` and runs `slurp-source` to copy it into `internal/config/config.go` (do not edit that file manually).
3. **`nix/openapi.nix`** — fetches the pfSense REST API OpenAPI JSON from GitHub releases and runs `patch-openapi` to flatten `allOf` entries, producing a spec compatible with the HashiCorp generator.
4. **`nix/config.nix`** — runs `gen-config` to produce a `config.yaml` (tfplugingen-openapi generator config) from the patched OpenAPI spec.
5. **`nix/provider-spec.nix`** — calls `a2b`'s `genProviderSpec` (wraps `tfplugingen-openapi`) with the config and OpenAPI spec to produce `schema.json`.
6. **`nix/provider.nix`** — calls `a2b`'s `genProvider` (wraps `tfplugingensdk`) with the schema to produce generated provider Go source. Compilation into a final provider binary is not yet complete.

### Go code roles

| Path                        | Purpose                                                                                                                                             |
| --------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------- |
| `cmd/patch-openapi/`        | CLI entry point for `PatchSpec`                                                                                                                     |
| `cmd/gen-config/`           | CLI entry point for `GenerateConfig`                                                                                                                |
| `cmd/slurp-source/`         | CLI entry point for `ExtractConfig`                                                                                                                 |
| `pkg/openapi.go`            | Flattens `allOf` schemas in an OpenAPI document via `libopenapi`                                                                                    |
| `pkg/config.go`             | Builds a `config.Config` from the OpenAPI model and writes YAML                                                                                     |
| `pkg/hack.go`               | Extracts `parse.go` from the upstream HashiCorp repo using Go's AST                                                                                 |
| `internal/config/config.go` | **Generated** — copied from `hashicorp/terraform-plugin-codegen-openapi`; defines `Config`, `Resource`, `DataSource`, `OpenApiSpecLocation` structs |

### Key dependencies

- `github.com/pb33f/libopenapi` — OpenAPI v3 parsing and bundling
- `github.com/unmango/go` — `world` package (OS abstraction for testability), `cli` utilities
- `gomod2nix` — keeps `nix/gomod2nix.toml` in sync with `go.mod` for reproducible Nix builds
- `a2b` (flake input `UnstoppableMango/a2b`) — provides `genProviderSpec` and `genProvider` Nix functions that wrap the HashiCorp tfplugingen-openapi toolchain

### Adding a new resource

Edit `pkg/config.go` → `ConfigFor()` to add an entry to the `Resources` map with the appropriate OpenAPI path/method for each CRUD operation. The downstream Nix pipeline will pick up the change on the next `nix build`.
