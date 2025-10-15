.PHONY: build clean test lint vet fmt deps install run

# Build variables
BINARY_NAME=crush-md
BUILD_DIR=bin
CMD_DIR=cmd/crush-md

# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet

# Build the application
build:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

# Build for multiple platforms
build-all: build-linux build-darwin build-windows

build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)

build-darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(CMD_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(CMD_DIR)

build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)

# Install dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Install the binary
install: build
	$(GOCMD) install ./$(CMD_DIR)

# Run the application
run:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR) && ./$(BUILD_DIR)/$(BINARY_NAME)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Format code
fmt:
	$(GOFMT) -s -w .

# Vet code
vet:
	$(GOVET) ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Check code quality
check: fmt vet test

# Development setup
dev-setup: deps
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Quick development build and test
dev: fmt vet test build

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  build-all      - Build for all platforms"
	@echo "  deps           - Download and tidy dependencies"
	@echo "  install        - Install the binary"
	@echo "  run            - Build and run the application"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  fmt            - Format code"
	@echo "  vet            - Vet code"
	@echo "  lint           - Lint code (requires golangci-lint)"
	@echo "  check          - Run format, vet, and test"
	@echo "  dev-setup      - Setup development environment"
	@echo "  dev            - Quick development build and test"
	@echo "  help           - Show this help"