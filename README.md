# Go Build Tools

`go-build-tools` is a Go-based Makefile utility designed to automate the development environment setup for projects.

Its main goal is to automate the installation of common tools used in Go development, such as protocol buffers, linting 
and mocks. It also provides a simple set of targets to cover the most common build, test and release scenarios.
It can be used locally or as part of a CI pipeline.

It also includes a set of GitHub workflows to automate the release process.

## Features

- **Linting**: Runs `golangci-lint` on all Go source files using a strict configuration file.
- **Testing**: Runs all Go tests with coverage.
- **Release**: Creates a GitHub release with the compiled binaries.
- **Protoc Generation**: Fetches the official `protoc` compiler release from GitHub and compiles your protocol buffers.
- **Mock Generation**: Generates mocks using `gomock` and `gomockhandler` for Go interfaces.
- **Code Formatting**: Runs `gofmt` on all Go source files.
- **CI Workflows**: GitHub Actions workflows to automate the release process that can be found in the `.github/workflows` directory.
They can be imported into your own GitHub workflows.

## Prerequisites

- **Go**: You need at least Go 1.24.
- **Make**: Make is required to run the Makefile.

## Installation

You can install these tools by cloning the repository as submodule in your project:

```bash
git submodule add https://github.com/Drathveloper/go-build-tools .parent
```

Then you can include the Makefile in your project:

```Makefile
PARENT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

include $(PARENT_DIR).parent/parent.mk

# Override variables here
SOURCE.DIR   := ./
```

## Variables:

The following variables can be overridden in your Makefile:

### Directories
- SRC.DIR: The root directory containing the source code.
- BIN.DIR: The directory where the compiled binaries will be placed.
- REPORTS.DIR: The directory where the test reports will be placed.
- LINT.CONFIG: The path to the linting configuration file.
- PB.DIR: The directory where the protocol buffers are located.
- SERVICE_NAME: The name of the service.

### Tools Versions
GOMOCK.VERSION: Gomock version. Default: 0.6.0
GOMOCKHANDLER.VERSION: Gomockhandler version. Default: 1.6.1
GOLANGCI_LINT.VERSION: GolangCI lint version. Default: v2.6.0
PROTOC.VERSION: Protoc compiler version. Default: 33.0
PROTOC_GEN.VERSION: Protoc gen plugin version. Default: v1.36.10
PROTOC_GEN_GRPC.VERSION: Protoc gen gRPC plugin version. Default: v1.5.1
