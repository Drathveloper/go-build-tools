GOMOCK.VERSION        := 0.6.0
GOMOCKHANDLER.VERSION := 1.6.1

GOMOCKHANDLER.BIN := gomockhandler

.PHONY: setup-mock-tools clean-mock-tools generate-mocks check-mocks

setup-mock-tools:
	@$(call log.info, Setup mock tools started)
	$(GO.BIN) install go.uber.org/mock/mockgen@v$(GOMOCK.VERSION)
	$(GO.BIN) install github.com/sanposhiho/gomockhandler@v$(GOMOCKHANDLER.VERSION)
	@$(call log.info, Setup mock tools finished successfully)

clean-mock-tools:
	@$(call log.info, Clean mock tools started)
	@$(call log.info, Clean mock tools finished successfully)

generate-mocks: setup-mock-tools
	@$(call log.info, Generate mocks started)
	@PATH=$$PATH:$(GO.DIR) $(GOMOCKHANDLER.BIN) mockgen || ( $(call log.error, Generate mocks failed) && false )
	@$(call log.info, Generate mocks finished successfully)

check-mocks: setup-mock-tools
	@$(call log.info, Check mocks started)
	@PATH=$$PATH:$(GO.DIR) $(GOMOCKHANDLER.BIN) check || ( $(call log.error, Check mocks failed) && false )
	@$(call log.info, Check mocks finished successfully)
