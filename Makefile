_ != mkdir -p bin

GOMOD2NIX ?= gomod2nix
TF_GEN ?= go tool tfplugingen-openapi

GO_SRC := $(shell find . -type f -name '*.go')

generate gen: bin/provider_code_spec.json

tidy: go.sum nix/gomod2nix.toml

go.sum: go.mod ${GO_SRC}
	go mod tidy

bin/provider_code_spec.json: bin/openapi.json generator_config.yml
	$(TF_GEN) generate $< \
		--config ./generator_config.yml \
		--output ./$@

bin/openapi.json:
	curl --fail https://pfrest.org/api-docs/openapi.json | jq -r >$@

nix/gomod2nix.toml: go.sum
	$(GOMOD2NIX) generate --outdir ./nix
