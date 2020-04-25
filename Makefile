SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif
.RECIPEPREFIX = >

export GO111MODULE := on
ROOTDIR := $(shell pwd)
VENDORDIR := $(ROOTDIR)/vendor
QUIET=@
VERIFYARGS ?=

GOOS ?=
GOOS := $(if $(GOOS),$(GOOS),linux)
GOARCH ?=
GOARCH := $(if $(GOARCH),$(GOARCH),amd64)
GOENV  := CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH)
GO     := $(GOENV) go
GO_BUILD := $(GO) build -trimpath

pkgs  = $(shell $(GO) list ./... | grep -v vendor)

help:
> @echo "Xene - Makefile"
> @echo ""
> @echo "Available make targets:"
> @echo ""
> @echo "* build: Run build script to build xene binary."
> @echo "* check-lint: Check if the formatting is correct for the go code, using golangci lint."
> @echo "* fix-lint: Fixes lint errors generated by golangci lint"
> @echo "* format: Format go code in the repository, run this before committing anything."
> @echo "* govet: Run govet on the code to check for any mistakes."
> @echo "* docs: Build docs site using mkdocs in site/ directory"
> @echo "* check-api-docs: checks the integrity of API docs for xene."
> @echo "* proto: Generate protobuf client and server code for definitions in pkg/proto/."
> @echo ""

# Build status
build:
> @echo "Building xene..."
> @./contrib/scripts/build.sh

# Check go formatting
check-lint:
> @echo "[*] Checking for formatting and linting errors"
> @./contrib/scripts/check-fmt.sh
> @golangci-lint run ./...

fix-lint:
> @echo "[*] Fixing lint errors using golangci-lint"
> @golangci-lint run ./... --fix

# Format code using gofmt
format:
> @echo "[*] Formatting code"
> @$(GO) fmt $(pkgs)

# Vet code using go vet
govet:
> @echo "[*] Vetting code, checking for mistakes"
> @$(GO) vet $(pkgs)

docs:
> @echo "[*] Building docs.."
> @rm -rf site/
> @mkdocs build
> @python contrib/scripts/swagger-docs.py

check-api-docs:
> @echo "[*] Checking API docs integrity"
> @./contrib/scripts/check-api-docs.sh

proto:
> @echo "[*] Generating proto definitions"
> @protoc -I pkg/proto pkg/proto/agent.proto --go_out=plugins=grpc:pkg/proto

.PHONY: build format check-lint fix-lint govet help docs proto
