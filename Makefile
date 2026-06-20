_ != mkdir -p bin

GOMOD2NIX ?= gomod2nix
NIX_SRC := $(shell find . -name '*.nix')

build:
	nix build .#

tools: nix/tools

update:
	nix flake update

check:
	nix flake check

nix/go.mod.patch: ${NIX_SRC}
	nix run .#bin.src.goModPatch -- $@

go.mod: nix/go.mod.patch
	nix build .#bin.src
	@cp result/$@ ${CURDIR}/$@

.PHONY: nix/tools
nix/tools:
	$(MAKE) -C $@
