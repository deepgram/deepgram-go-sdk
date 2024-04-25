# Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
# Use of this source code is governed by a MIT license that can be found in the LICENSE file.
# SPDX-License-Identifier: MIT

# detect the build OS
ifeq ($(OS),Windows_NT)
	build_OS := Windows
	NUL = NUL
else
	build_OS := $(shell uname -s 2>/dev/null || echo Unknown)
	NUL = /dev/null
endif

.DEFAULT_GOAL:=help

##### GLOBAL
ROOT_DIR := $(shell git rev-parse --show-toplevel)
GO := go
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOHOSTOS ?= $(shell go env GOHOSTOS)
GOHOSTARCH ?= $(shell go env GOHOSTARCH)

# Add tooling binaries here and in hack/tools/Makefile
TOOLS_BIN_DIR := $(shell mktemp -d)
GOLANGCI_LINT := $(TOOLS_BIN_DIR)/golangci-lint

help: #### Display help
	@echo ''
	@echo 'Syntax: make <target>'
	@awk 'BEGIN {FS = ":.*## "; printf "\nTargets:\n"} /^[a-zA-Z_-]+:.*?#### / { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ''
##### GLOBAL

##### BUILD
export GOHOSTOS
export GOHOSTARCH
export GO
export GOLANGCI_LINT
##### BUILD

##### LINTING TARGETS
.PHONY: version
version: #### display version of components
	@echo 'ROOT_DIR: $(ROOT_DIR)'
	@echo 'GOOS: $(GOOS)'
	@echo 'GOARCH: $(GOARCH)'
	@echo 'go version: $(shell go version)'

.PHONY: lint mdlint shellcheck actionlint yamllint
lint:
	@echo 'Noop'

mdlint:
	@echo 'Noop'

shellcheck:
	@echo 'Noop'

actionlint:
	@echo 'Noop'

yamllint:
	@echo 'Noop'

