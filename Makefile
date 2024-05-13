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

.PHONY: check lint mdlint shellcheck actionlint yamllint ### Performs all of the checks, lint'ing, etc available
check: lint mdlint shellcheck actionlint yamllint

.PHONY: ensure-deps
ensure-deps: #### Ensure that all required dependency utilities are downloaded or installed
	hack/ensure-deps/ensure-dependencies.sh

GO_MODULES=$(shell find . -path "*/go.mod" | xargs -I _ dirname _)

# get-deps:
# 	@for i in $(GO_MODULES); do \
# 		echo "-- Getting deps for $$i --"; \
# 		working_dir=`pwd`; \
# 		if [ "$${i}" = "." ]; then \
# 			go mod tidy; \
# 		else \
# 			cd $${i}; \
# 			$(MAKE) get-deps || exit 1; \
# 			cd $$working_dir; \
# 		fi; \
# 	done

# verify-modules: get-deps
# 	@for i in $(GO_MODULES); do \
# 		echo "-- Verifying modules for $$i --"; \
# 		working_dir=`pwd`; \
# 		cd $${i}; \
# 		if [ "`git diff --name-only HEAD -- go.sum go.mod`" != "" ]; then \
# 			echo "go module files in $$i directory are out of date, run 'go mod tidy' and commit the changes"; exit 1; \
# 		fi; \
# 		cd $$working_dir; \
# 	done

lint: tools #### Performs Golang programming lint
	@for i in $(GO_MODULES); do \
		echo "-- Linting $$i --"; \
		working_dir=`pwd`; \
		if [ "$${i}" = "." ]; then \
			$(GOLANGCI_LINT) run -v --timeout=5m; \
		else \
			cd $${i}; \
			$(MAKE) lint || exit 1; \
			cd $$working_dir; \
		fi; \
	done; \
	CHECK=$$(grep -r --include="*.go" ioutil ./); \
	if [ -n "$${CHECK}" ]; then \
		echo "ioutil is deprecated, use io or os replacements"; \
		echo "$${CHECK}"; \
		exit 1; \
	fi

mdlint: #### Performs Markdown lint
	# mdlint rules with common errors and possible fixes can be found here:
	# https://github.com/DavidAnson/markdownlint/blob/main/doc/Rules.md
	hack/check/check-mdlint.sh

shellcheck: #### Performs bash/shell lint
	hack/check/check-shell.sh

yamllint: #### Performs yaml lint
	hack/check/check-yaml.sh

actionlint: #### Performs GitHub Actions lint
	actionlint
##### LINTING TARGETS

##### Tooling Binaries
TOOLS_DIR := $(ROOT_DIR)/hack/check/tools

tools: $(GOLANGCI_LINT)
.PHONY: $(GOLANGCI_LINT)
$(GOLANGCI_LINT):
	make -C $(TOOLS_DIR) $(@F)
##### Tooling Binaries
