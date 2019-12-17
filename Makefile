GO := go
pkgs  = $(shell $(GO) list ./... | grep -v vendor)

help:
	@echo "Xene - Makefile"

# Build status
build:
	@./scripts/build/build.sh

# Check go formatting
check_format:
	@echo "[*] Checking for formatting errors using gofmt"
	@./scripts/build/check_gofmt.sh

# Format code using gofmt
format:
	@echo "[*] Formatting code"
	@$(GO) fmt $(pkgs)

# Vet code using go vet
govet:
	@echo "[*] Vetting code, checking for mistakes"
	@$(GO) vet $(pkgs)

.PHONY: build format check_format govet help

