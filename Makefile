# Makefile for Cycles CPU Monitor
# Provides convenient build, test, and development commands

.PHONY: all build run test clean install-deps fmt vet help

# Variables
BINARY_NAME=cycles
BUILD_DIR=build
VERSION=$(shell grep 'Version:' config.go | cut -d'"' -f2)
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION)"

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

# Default target
all: build

## help: Display this help message
help:
	@echo "Cycles Build System"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' Makefile | sed 's/## /  /'

## build: Build a single bundled binary
build:
	@echo "$(GREEN)Building Cycles v$(VERSION)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "$(GREEN)✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

## run: Build and run the application
run: build
	@echo "$(GREEN)Running Cycles...$(NC)"
	@./$(BUILD_DIR)/$(BINARY_NAME)

## test: Run all tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	@go test -v ./...

## test-short: Run tests without verbose output
test-short:
	@echo "$(GREEN)Running tests...$(NC)"
	@go test ./...

## bench: Run benchmarks
bench:
	@echo "$(GREEN)Running benchmarks...$(NC)"
	@go test -bench=. -benchmem ./...

## fmt: Format all Go source files
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

## vet: Run go vet on all source files
vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ Vet complete$(NC)"

## check: Run fmt, vet, and test
check: fmt vet test

## clean: Remove build artifacts
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)
	@echo "$(GREEN)✓ Clean complete$(NC)"

## install-deps: Install/update Go dependencies
install-deps:
	@echo "$(GREEN)Installing dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✓ Dependencies installed$(NC)"

## setup: Run development environment setup
setup:
	@echo "$(GREEN)Running development setup...$(NC)"
	@chmod +x scripts/setup-dev.sh
	@./scripts/setup-dev.sh

## release: Build optimized release binary
release:
	@echo "$(GREEN)Building release binary v$(VERSION)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=1 go build $(LDFLAGS) -trimpath -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "$(GREEN)✓ Release build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

## dev: Quick development build (no optimization)
dev:
	@echo "$(GREEN)Building development binary...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "$(GREEN)✓ Dev build complete$(NC)"

## watch: Watch for changes and rebuild (requires entr)
watch:
	@echo "$(GREEN)Watching for changes...$(NC)"
	@echo "$(YELLOW)Note: Requires 'entr' (install: apt-get install entr)$(NC)"
	@find . -name '*.go' | entr -r make dev run

## size: Show binary size information
size: build
	@echo "$(GREEN)Binary size information:$(NC)"
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME)
	@file $(BUILD_DIR)/$(BINARY_NAME)

## install: Install the binary to /usr/local/bin (requires sudo)
install: release
	@echo "$(GREEN)Installing to /usr/local/bin...$(NC)"
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)✓ Installed successfully$(NC)"

## uninstall: Remove the installed binary
uninstall:
	@echo "$(YELLOW)Removing from /usr/local/bin...$(NC)"
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)✓ Uninstalled successfully$(NC)"

## docker-build: Build using Docker (no system dependencies needed)
docker-build:
	@echo "$(GREEN)Building with Docker...$(NC)"
	@docker build -t cycles-builder -f Dockerfile.build .
	@docker run --rm -v $(PWD)/$(BUILD_DIR):/output cycles-builder
	@echo "$(GREEN)✓ Docker build complete$(NC)"

## info: Display build environment information
info:
	@echo "$(GREEN)Build Environment Information:$(NC)"
	@echo "  Version:     $(VERSION)"
	@echo "  Go Version:  $(shell go version)"
	@echo "  Build Dir:   $(BUILD_DIR)"
	@echo "  Binary Name: $(BINARY_NAME)"
	@echo "  GOOS:        $(shell go env GOOS)"
	@echo "  GOARCH:      $(shell go env GOARCH)"
