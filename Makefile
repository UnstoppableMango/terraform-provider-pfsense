_ != mkdir -p bin

GOMOD2NIX ?= gomod2nix
TF_GEN ?= go tool tfplugingen-openapi

GO_SRC := $(shell find . -type f -name '*.go')

build: tidy internal/config/config.go
	nix build .#binary

.PHONY: provider-deps
provider-deps:
	go get github.com/hashicorp/terraform-plugin-framework@latest
	go get github.com/hashicorp/terraform-plugin-framework-validators@latest
	$(MAKE) tidy

update:
	nix flake update

check:
	nix flake check

tidy: go.sum nix/gomod2nix.toml

go.sum: go.mod ${GO_SRC}
	go mod tidy

.SECONDARY: internal/config/config.go
internal/config/config.go: nix/upstream.nix
	@mkdir -p $(@D) && rm -f $@
	cp $$(nix build .#upstream --print-out-paths) $@

nix/gomod2nix.toml: go.sum
	$(GOMOD2NIX) generate --outdir ./nix
