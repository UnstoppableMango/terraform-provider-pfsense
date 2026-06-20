GO_SRC := $(shell find . -type f -name '*.go')

tidy: go.sum gomod2nix.toml

go.sum: go.mod ${GO_SRC}
	go mod tidy

.SECONDARY: internal/config/config.go
internal/config/config.go: nix/upstream.nix
	@mkdir -p $(@D) && rm -f $@
	cp $$(nix build ../..#upstream --print-out-paths) $@

gomod2nix.toml: go.sum
	$(GOMOD2NIX) generate
