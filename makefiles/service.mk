.PHONY: fmt clean run build

SERVICE_NAME := service

SOURCE.DIR   := ./src
REPORTS.DIR  := ./reports

SOURCE.PKGS  := $(shell $(GO.BIN) list $(SOURCE.DIR)/... | grep -v "vendor" 2>/dev/null)

fmt:
	@$(call log.info, Code format started)
	@$(GO.BIN) fmt $(SOURCE.DIR) || ( $(call log.error, Code format failed) && false )
	@$(call log.info, Code formatted successfully)

clean:
	@$(call log.info, Clean binary executables started)
	@rm -rf bin/* >/dev/null 2>&1 || true
	@$(call log.info, Binary executables cleaned successfully)

run:
	@$(call log.info, Run go application started)
	@$(GO.BIN) run $(SOURCE.DIR) || ( $(call log.error, Run go application failed) && false )
	@$(call log.info, Run go application finished)

build:
	@$(call log.info, Build binary executable started)
	@$(GO.BIN) build -o bin/$(SERVICE_NAME) $(SOURCE.DIR) || ( $(call log.error, Build binary executable failed) && false )
	@$(call log.info, Binary executable builded successfully)

test:
	@$(call log.info, Run unit tests started)
	@mkdir -p $(REPORTS.DIR)
	$(GO.BIN) test -v $(SOURCE.PKGS) -coverprofile $(REPORTS.DIR)/coverage || ( $(call log.error, Run unit tests failed) && false )
	@$(call log.info, Unit tests finished successfully)
