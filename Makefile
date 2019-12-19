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


GO := go
pkgs  = $(shell $(GO) list ./... | grep -v vendor)

help:
> @echo "Xene - Makefile"
> @echo ""
> @echo "Available make targets:"
> @echo ""
> @echo "* build: Run build script to build xene binary."
> @echo "* check_format: Check if the formatting is correct for the go code, using gofmt."
> @echo "* format: Format go code in the repository, run this before committing anything."
> @echo "* govet: Run govet on the code to check for any mistakes."

# Build status
build:
> @echo "Building xene..."
> @./scripts/build/build.sh

# Check go formatting
check_format:
> @echo "[*] Checking for formatting errors using gofmt"
> @./scripts/build/check_gofmt.sh

# Format code using gofmt
format:
> @echo "[*] Formatting code"
> @$(GO) fmt $(pkgs)

# Vet code using go vet
govet:
> @echo "[*] Vetting code, checking for mistakes"
> @$(GO) vet $(pkgs)

.PHONY: build format check_format govet help
