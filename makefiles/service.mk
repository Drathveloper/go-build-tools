fmt:
	@$(GO.BIN) fmt ./...
	@$(call log.info, Code formatted successfully)

clean:
	@rm -rf bin/* >/dev/null 2>&1 || true
	@$(call log.info, Binaries cleaned successfully)

run:
	$(GO.BIN) run ./...

build:
	$(GO.BIN) build
