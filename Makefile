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

generate gen: schema.json nix/go.mod.patch nix/provider-gomod2nix.toml

schema.json: ${GEN_SRC} ${NIX_SRC} flake.lock
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

nix/go.mod.patch: schema.json ${NIX_SRC} flake.lock
	nix run .#bin.src.goModPatch -- $@

nix/provider-gomod2nix.toml: nix/go.mod.patch
	nix run .#bin.src.gomod2nixToml -- $@

nix/gomod2nix.toml: go.sum
	$(GOMOD2NIX) generate --outdir ./nix
