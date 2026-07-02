_ != mkdir -p bin

GOMOD2NIX ?= gomod2nix
NIX_SRC := $(shell find . -name '*.nix')

build: generate
	nix build .#

generate gen: nix/go.mod.patch nix/gomod2nix.toml.patch

src:
	nix build .#bin.src

tools:
	nix build .#tools

update:
	nix flake update

check: generate
	nix flake check

tidy:
	$(MAKE) -C nix/tools tidy

nix/go.mod.patch: ${NIX_SRC} flake.lock
	nix run .#bin.src.goModPatch -- $@

nix/gomod2nix.toml.patch: nix/go.mod.patch
	nix run .#bin.src.gomod2nixTomlPatch -- $@

nix/gomod2nix.toml: nix/go.mod.patch
	nix run .#bin.src.gomod2nixToml -- ${@D}

