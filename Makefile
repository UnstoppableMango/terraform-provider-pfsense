_ != mkdir -p bin

GOMOD2NIX ?= gomod2nix
TF_GEN ?= go tool tfplugingen-openapi

build: tidy
	nix build .#

update:
	nix flake update

check:
	nix flake check
