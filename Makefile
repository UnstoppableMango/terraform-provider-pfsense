_ != mkdir -p bin

GOMOD2NIX ?= gomod2nix
NIX_SRC := $(shell find . -name '*.nix')

build: generate
	nix build .#

.PHONY: test
test:
	nix run .#tests

generate gen: nix/go.mod.patch nix/gomod2nix.toml.patch

src:
	nix build .#bin.src

tidy: go.sum nix/gomod2nix.toml

tools:
	nix build .#tools

update:
	nix flake update

check: generate
	nix flake check

go.sum: go.mod
	go mod tidy

nix/go.mod.patch: ${NIX_SRC} flake.lock
	nix run .#bin.src.goModPatch -- $@

nix/gomod2nix.toml.patch: nix/go.mod.patch
	nix run .#bin.src.gomod2nixTomlPatch -- $@

nix/gomod2nix.toml: go.sum
	$(GOMOD2NIX) generate --outdir ./nix
