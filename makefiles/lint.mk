LINTER.BIN := golangci-lint

THIS.MAKEFILE := $(lastword $(MAKEFILE_LIST))
THIS.DIR := $(dir $(abspath $(THIS.MAKEFILE)))

.PHONY: setup-lint lint lint-fix lint-clean

setup-lint:
	@$(call log.info, Setup golangci-lint started)
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GO.DIR) $(GOLANGCI_LINT.VERSION)
	@$(call log.info, Setup golangci-lint finished successfully)

lint: setup-lint
	@$(call log.info, Lint started)
	@PATH=$$PATH:$(GO.DIR) $(LINTER.BIN) --config $(THIS.DIR)../configs/golangci.yml run || ( $(call log.error, Lint failed) && false )
	@$(call log.info, Lint finished successfully)

lint-fix: setup-lint
	@$(call log.info, Lint with fix started)
	@PATH=$$PATH:$(GO.DIR) $(LINTER.BIN) --config $(THIS.DIR)../configs/golangci.yml run --fix || ( $(call log.error, Lint failed) && false )
	@$(call log.info, Lint with fix finished successfully)

lint-clean:
	@$(call log.info, Clean lint cache started)
	@PATH=$$PATH:$(GO.DIR) $(LINTER.BIN) cache clean
	@$(call log.info, Clean lint cache finished successfully)
