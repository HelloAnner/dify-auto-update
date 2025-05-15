 # Go parameters
BINARY_NAME=dify-auto-update
MAIN_PATH=cmd/main.go
GO=go

# Build directory
BUILD_DIR=build

# Version and commit information
VERSION?=1.0.0
COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

# Supported platforms
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

# Clean build directory
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# Create build directory
.PHONY: init
init:
	mkdir -p $(BUILD_DIR)

# Build for current platform
.PHONY: build
build:
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Build for all platforms
.PHONY: build-all
build-all: clean init
	$(foreach platform,$(PLATFORMS),\
		$(eval GOOS=$(word 1,$(subst /, ,$(platform))))\
		$(eval GOARCH=$(word 2,$(subst /, ,$(platform))))\
		$(eval EXTENSION=$(if $(filter windows,$(GOOS)),.exe))\
		$(eval OUTFILE=$(BUILD_DIR)/$(BINARY_NAME)_$(GOOS)_$(GOARCH)$(EXTENSION))\
		GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build $(LDFLAGS) -o $(OUTFILE) $(MAIN_PATH);\
	)

# Build for specific platforms
.PHONY: build-darwin-amd64
build-darwin-amd64: init
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)_darwin_amd64 $(MAIN_PATH)

.PHONY: build-darwin-arm64
build-darwin-arm64: init
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)_darwin_arm64 $(MAIN_PATH)

.PHONY: build-linux-amd64
build-linux-amd64: init
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)_linux_amd64 $(MAIN_PATH)

.PHONY: build-linux-arm64
build-linux-arm64: init
	GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)_linux_arm64 $(MAIN_PATH)

.PHONY: build-windows-amd64
build-windows-amd64: init
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)_windows_amd64.exe $(MAIN_PATH)

# Run tests
.PHONY: test
test:
	$(GO) test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

# Install dependencies
.PHONY: deps
deps:
	$(GO) mod download
	$(GO) mod verify

# Format code
.PHONY: fmt
fmt:
	$(GO) fmt ./...

# Run linter
.PHONY: lint
lint:
	$(if $(shell command -v golangci-lint 2> /dev/null),\
		golangci-lint run,\
		$(error "golangci-lint is not installed. Please install it first."))

# Default target
.DEFAULT_GOAL := build