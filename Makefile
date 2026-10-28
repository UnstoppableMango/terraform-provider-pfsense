_ != mkdir -p bin

GOMOD2NIX ?= gomod2nix
GO_SRC := $(shell find . -name '*.go')
NIX_SRC := $(shell find . -name '*.nix')
GEN_SRC := $(shell find cmd internal -name '*.go' -not -name '*_test.go')

build: generate
	nix build .#

.PHONY: test
test:
	nix run .#tests

.PHONY: generate gen
generate gen: gen/schema.json gen/go.mod.patch gen/gomod2nix.toml

gen/schema.json: ${GEN_SRC} ${NIX_SRC} flake.lock go.mod go.sum nix/gomod2nix.toml
	out=$$(nix build .#bin.spec --no-link --print-out-paths) && install -m 644 "$$out" $@

src:
	nix build .#bin.src

tidy: go.sum nix/gomod2nix.toml

tools:
	nix build .#tools

update:
	nix flake update

check: generate
	nix flake check

go.sum: go.mod ${GO_SRC}
	go mod tidy

gen/go.mod.patch: gen/schema.json ${NIX_SRC} flake.lock
	nix run .#bin.src.goModPatch -- $@

gen/gomod2nix.toml: gen/go.mod.patch
	nix run .#bin.src.gomod2nixToml -- $@

nix/gomod2nix.toml: go.sum
	$(GOMOD2NIX) generate --outdir ./nix
