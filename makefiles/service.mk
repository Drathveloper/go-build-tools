SERVICE_NAME := service

SOURCE.DIR   ?= ./src
BIN.DIR      ?= ./bin
REPORTS.DIR  ?= ./reports

SOURCE.PKGS  = $(shell $(GO.BIN) list $(SOURCE.DIR)/... | grep -v "vendor" 2>/dev/null)

.PHONY: download-deps fmt clean run compile build test benchmark

download-deps:
	@$(call log.info, Download go.mod dependencies started)
	@$(GO.BIN) mod download || ( $(call log.error, Code format failed) && false )
	@$(call log.info, go.mod dependencies downloaded successfully)

fmt:
	@$(call log.info, Code format started)
	@$(GO.BIN) fmt $(SOURCE.DIR) || ( $(call log.error, Code format failed) && false )
	@$(call log.info, Code formatted successfully)

clean:
	@$(call log.info, Clean binary executables started)
	@rm -rf $(BIN.DIR) >/dev/null 2>&1 || true
	@rm -rf $(REPORTS.DIR) >/dev/null 2>&1 || true
	@$(call log.info, Binary executables cleaned successfully)

run:
	@$(call log.info, Run go application started)
	@$(GO.BIN) run $(SOURCE.DIR) || ( $(call log.error, Run go application failed) && false )
	@$(call log.info, Run go application finished)

compile: generate-grpc
	@$(call log.info, Compile go application started)
	@$(GO.BIN) build $(SOURCE.DIR) || ( $(call log.error, Run go application failed) && false )
	@$(call log.info, Compile go application finished)

build: generate-grpc
	@$(call log.info, Build binary executable started)
	@$(GO.BIN) build -o $(BIN.DIR)/$(SERVICE_NAME) $(SOURCE.DIR) || ( $(call log.error, Build binary executable failed) && false )
	@$(call log.info, Binary executable builded successfully)

test: check-mocks
	@$(call log.info, Run unit tests started)
	@mkdir -p $(REPORTS.DIR)
	@$(GO.BIN) test -v $(SOURCE.PKGS) -coverprofile $(REPORTS.DIR)/coverage || ( $(call log.error, Run unit tests failed) && false )
	@$(call log.info, Unit tests finished successfully)

benchmark:
	@$(call log.info, Run benchmark tests started)
	@mkdir -p $(REPORTS.DIR)
	@$(GO.BIN) test -bench=. -benchmem -memprofile=$(REPORTS.DIR)/mem.out -cpuprofile=$(REPORTS.DIR)/cpu.out > $(REPORTS.DIR)/benchmarks || ( $(call log.error, Run benchmark tests failed) && false )
	@$(call log.info, Benchmark tests finished successfully)