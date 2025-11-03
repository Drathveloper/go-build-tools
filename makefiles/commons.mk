GO.BIN     = go
GO.ARCH    = $(shell $(GO.BIN) env GOARCH 2>/dev/null)
GO.OS      = $(shell $(GO.BIN) env GOOS 2>/dev/null)
GO.PATH    = $(shell $(GO.BIN) env GOPATH 2>/dev/null)
GO.VERSION = $(shell $(GO.BIN) env GOVERSION 2>/dev/null)
GO.DIR     = $(shell $(GO.BIN) env GOBIN)
ifeq ($(GO.DIR),)
GO.DIR     = $(shell $(GO.BIN) env GOPATH)/bin
endif

CURRENT.DIR := $(dir $(abspath $(firstword $(MAKEFILE_LIST))))

TOOLS.DIR   := tools

Logs.Color.NoColor=\033[0m
Logs.Color.Black=\033[1;30m
Logs.Color.Red=\033[1;31m
Logs.Color.Green=\033[1;32m
Logs.Color.Yellow=\033[1;33m
Logs.Color.White=\033[1;37m

define log.info
	printf "$(BC.White)%s $(Logs.Color.Green)INFO$(Logs.Color.NoColor)%s\n" "$(shell date '+%Y-%m-%dT%H:%M:%SZ')" "$(1)"
endef

define log.warn
	printf "$(BC.White)%s $(Logs.Color.Yellow)WARN$(Logs.Color.NoColor)%s\n" "$(shell date '+%Y-%m-%dT%H:%M:%SZ')" "$(1)"
endef

define log.error
	printf "$(BC.White)%s $(Logs.Color.Red)ERROR$(Logs.Color.NoColor)%s\n" "$(shell date '+%Y-%m-%dT%H:%M:%SZ')" "$(1)"
endef
