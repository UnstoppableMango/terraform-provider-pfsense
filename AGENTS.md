# AGENTS.md

This file provides guidance to AI agents when working with code in this repository.

## Commands

```bash
# Build (uses nix)
make build          # or: nix build

# Run tests
cd nix/tools && go test ./...

# Lint / format
nix flake check     # runs actionlint, nixfmt, gofmt via treefmt

# Update go dependencies
make tidy   # runs go mod tidy + gomod2nix in nix/tools/

# Update nix flake inputs
make update         # or: nix flake update

# Build individual nix packages
nix build .#tools       # CLI tools binary
nix build .#bin.src     # generated provider Go source
nix build               # final provider binary (default)
```

## Architecture

This project generates a Terraform provider for pfSense entirely from the pfSense REST API's OpenAPI spec, using a Nix-driven pipeline. No provider code is written by hand — it is all generated.

### Code generation pipeline (nix build order)

1. **`nix/tools/`** — separate Go submodule; builds CLI tools (`patch-openapi`, `gen-config`, `slurp-source`, `patch-provider`) and fetches `config.go` from upstream HashiCorp repo via `slurp-source` (do not edit `nix/tools/internal/config/config.go` manually).
2. **`nix/openapi.nix`** — fetches pfSense REST API OpenAPI JSON from GitHub releases; runs `patch-openapi` to flatten `allOf` entries, producing a spec compatible with the HashiCorp generator.
3. **`nix/provider-spec.nix`** — runs `gen-config` then calls `a2b`'s `genProviderSpec` (wraps `tfplugingen-openapi`) to produce `schema.json`.
4. **`nix/provider-src.nix`** — calls `a2b`'s `genProvider` + `scaffold` (wraps `tfplugingensdk`) with the schema to produce generated provider Go source; `nix/default.nix` compiles it into the final provider binary.

### Go code roles

| Path                                      | Purpose                                                                                                                                                   |
| ----------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `cmd/terraform-provider-pfsense/`         | Provider binary entry point (thin `main.go`)                                                                                                              |
| `nix/tools/cmd/patch-openapi/`            | Flattens `allOf` schemas in OpenAPI doc via `libopenapi`                                                                                                  |
| `nix/tools/cmd/gen-config/`               | Builds `config.Config` from OpenAPI model and writes YAML                                                                                                 |
| `nix/tools/cmd/slurp-source/`             | Extracts `parse.go` from upstream HashiCorp repo using Go's AST                                                                                           |
| `nix/tools/cmd/patch-provider/`           | Patches the scaffolded provider Go source                                                                                                                 |
| `nix/tools/internal/config/config.go`     | **Generated** — copied from `hashicorp/terraform-plugin-codegen-openapi`; defines `Config`, `Resource`, `DataSource`, `OpenApiSpecLocation` structs       |

### Key dependencies

- `github.com/pb33f/libopenapi` — OpenAPI v3 parsing and bundling
- `github.com/unmango/go` — `world` package (OS abstraction for testability), `cli` utilities
- `gomod2nix` — keeps `nix/gomod2nix.toml` in sync with `go.mod` for reproducible Nix builds
- `a2b` (flake input `UnstoppableMango/a2b`) — provides `genProviderSpec` and `genProvider` Nix functions that wrap the HashiCorp tfplugingen-openapi toolchain

### Adding a new resource

Edit `nix/tools/cmd/gen-config/config.go` → `ConfigFor()` to add an entry to the `Resources` map with the appropriate OpenAPI path/method for each CRUD operation. The downstream Nix pipeline will pick up the change on the next `nix build`.
