_ != mkdir -p bin

TF_GEN ?= go tool tfplugingen-openapi

generate gen: bin/provider_code_spec.json

bin/provider_code_spec.json: bin/openapi.json generator_config.yml
	$(TF_GEN) generate $< \
		--config ./generator_config.yml \
		--output ./$@

bin/openapi.json:
	curl --fail https://pfrest.org/api-docs/openapi.json | jq -r >$@
