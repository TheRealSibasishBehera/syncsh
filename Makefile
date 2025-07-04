# Variables
BINARY_NAME=syncsh
GO_FILES=$(shell find . -name "*.go" -type f)
BUILD_DIR=build
PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: all build clean test fmt vet lint run help install cross-build

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) .

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@go clean

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Run the application
run:
	@go run .

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Cross-compile for multiple platforms
cross-build:
	@echo "Cross-compiling for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d'/' -f1); \
		arch=$$(echo $$platform | cut -d'/' -f2); \
		output=$(BUILD_DIR)/$(BINARY_NAME)-$$os-$$arch; \
		if [ $$os = "windows" ]; then output="$$output.exe"; fi; \
		echo "Building for $$os/$$arch..."; \
		GOOS=$$os GOARCH=$$arch go build -o $$output .; \
	done

# Help target
help:
	@echo "Available targets:"
	@echo "  all          - Build the binary (default)"
	@echo "  build        - Build the binary"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  fmt          - Format code"
	@echo "  vet          - Run go vet"
	@echo "  run          - Run the application"
	@echo "  deps         - Install and tidy dependencies"
	@echo "  install      - Install binary to /usr/local/bin"
	@echo "  cross-build  - Cross-compile for multiple platforms"
	@echo "  help         - Show this help message"