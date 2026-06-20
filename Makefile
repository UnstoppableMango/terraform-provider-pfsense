_ != mkdir -p bin

GOMOD2NIX ?= gomod2nix
NIX_SRC := $(shell find . -name '*.nix')

build: nix/gomod2nix.toml
	nix build .#

tools: nix/tools

update:
	nix flake update

check:
	nix flake check

nix/go.mod.patch: ${NIX_SRC}
	nix run .#bin.src.goModPatch -- $@

nix/gomod2nix.toml: nix/go.mod.patch
	nix run .#bin.src.gomod2nixToml -- ${@D}

go.mod go.sum &: nix/go.mod.patch
	nix build .#bin.src
	install -m 644 result/go.{mod,sum} ${CURDIR}/

.PHONY: nix/tools
nix/tools:
	$(MAKE) -C $@
