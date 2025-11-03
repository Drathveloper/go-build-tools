PB.DIR        := proto
PROTOC.BIN    := $(TOOLS.DIR)/protoc

.PHONY: setup-grpc-tools clean-grpc-tools generate-grpc

setup-grpc-tools:
	@$(call log.info, Setup gRPC tools started)
	$(GO.BIN) -C scripts run cmd/setup-grpc/main.go $(CURRENT.DIR) $(PROTOC.VERSION)
	$(GO.BIN) install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN.VERSION)
	$(GO.BIN) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GRPC.VERSION)
	@$(call log.info, Setup gRPC tools finished successfully)

clean-grpc-tools:
	@$(call log.info, Clean gRPC tools started)
	@rm -rf $(PROTOC.BIN)
	@rm -f "$(GO.DIR)/protoc-gen-go"
	@rm -f "$(GO.DIR)/protoc-gen-go-grpc"
	@$(call log.info, Clean gRPC tools finished successfully)

generate-grpc: setup-grpc-tools
	@$(call log.info, Generate gRPC stubs started)
ifneq ($(shell find $(PB.DIR) -name "*.proto" 2>/dev/null),)
	$(PROTOC.BIN) -I $(PB.DIR) --go_out=. --go-grpc_out=. $(shell find $(PB.DIR) -name "*.proto" 2>/dev/null)
endif
	@$(call log.info, Generate gRPC stubs finished successfully)
