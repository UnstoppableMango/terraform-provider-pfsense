_ != mkdir -p bin

GOMOD2NIX ?= gomod2nix
NIX_SRC := $(shell find . -name '*.nix')

build: nix/gomod2nix.toml
	nix build .#

src:
	nix build .#bin.src

tools:
	nix build .#tools

update:
	nix flake update

check:
	nix flake check

tidy:
	$(MAKE) -C nix/tools tidy

nix/go.mod.patch: ${NIX_SRC} flake.lock
	nix run .#bin.src.goModPatch -- $@

nix/gomod2nix.toml: nix/go.mod.patch
	nix run .#bin.src.gomod2nixToml -- ${@D}

go.mod go.sum &: nix/go.mod.patch
	nix build .#bin.src
	install -m 444 result/go.{mod,sum} ${CURDIR}/
