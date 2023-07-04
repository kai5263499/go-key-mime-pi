which = $(shell which $1 2> /dev/null || echo $1)

GO_PATH := $(call which,go)
$(GO_PATH):
	$(error Missing go)

MOQ_PATH := $(call which,moq)
$(MOQ_PATH):
	@$(GO_PATH) install github.com/matryer/moq@latest


LINTER_PATH := $(call which,golangci-lint)
$(LINTER_PATH):
	$(error Missing golangci: https://golangci-lint.run/usage/install)
lint:
	@rm -rf ./vendor
	@$(GO_PATH) mod vendor
	export GOMODCACHE=./vendor
	@$(LINTER_PATH) run

.PHONY: test
test:
	@$(GO_PATH) test -v -cover ./...

OAPI_PATH := $(call which,oapi-codegen)
$(OAPI_PATH):
	$(error Missing oapi-codegen: https://github.com/deepmap/oapi-codegen)

internal/%/server.gen.go: api/%.yaml 
	@$(OAPI_PATH) -package $(notdir $(@D)) -include-tags="$*" -generate spec,server $< > $@

internal/%/types.gen.go: api/%.yaml 
	@$(OAPI_PATH) -package $(notdir $(@D)) -generate types $< > $@

api/v1: $(OAPI_PATH)
	@$(GO_PATH) generate -v ./internal/v1/...
.PHONY: api/v1