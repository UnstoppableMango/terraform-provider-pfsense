GO_SRC := $(shell find ./cmd -type f -name '*.go')

tidy: go.sum gomod2nix.toml

go.sum: go.mod ${GO_SRC}
	go mod tidy

.SECONDARY: internal/config/config.go
internal/config/config.go: default.nix
	@mkdir -p $(@D) && rm -f $@
	cp $$(nix build ../..#tools.configGo --print-out-paths) $@

gomod2nix.toml: go.sum
	$(GOMOD2NIX) generate
