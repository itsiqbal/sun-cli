# Makefile for mycli - Optimized Go CLI build

# Metadata
APP_NAME := sun
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := $(GOCMD) fmt
GOVET := $(GOCMD) vet

# Build parameters
BINARY_NAME := $(APP_NAME)
BINARY_UNIX := $(BINARY_NAME)_unix
BINARY_WINDOWS := $(BINARY_NAME).exe
BINARY_MAC := $(BINARY_NAME)_darwin

# Directories
BUILD_DIR := build
DIST_DIR := dist
CMD_DIR := cmd

# Ldflags for version injection
LDFLAGS := -ldflags "\
	-s -w \
	-X main.version=$(VERSION) \
	-X main.commit=$(COMMIT) \
	-X main.date=$(BUILD_DATE)"

# Optimization flags
GCFLAGS := -gcflags="all=-trimpath=$(PWD)"
ASMFLAGS := -asmflags="all=-trimpath=$(PWD)"

# Build tags for reduced binary size
BUILD_TAGS := -tags="netgo,osusergo"

# Platform-specific settings
PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

.PHONY: all build clean test coverage lint fmt vet deps help install uninstall run
.PHONY: build-all build-linux build-mac build-windows release docker
.DEFAULT_GOAL := help

## help: Display this help message
help:
	@echo "$(APP_NAME) - Optimized Build System"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_-]+:.*?##/ { printf "  %-15s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

## build: Build the binary for current platform (optimized)
build: deps
	@echo "🔨 Building $(BINARY_NAME) v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) $(GCFLAGS) $(ASMFLAGS) $(BUILD_TAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME)

## build-fast: Fast build without optimizations (for development)
build-fast:
	@echo "⚡ Fast building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Fast build complete"

## build-all: Build for all platforms
build-all: clean
	@echo "🌍 Building for all platforms..."
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		OS=$$(echo $$platform | cut -d'/' -f1); \
		ARCH=$$(echo $$platform | cut -d'/' -f2); \
		OUTPUT=$(DIST_DIR)/$(BINARY_NAME)-$$OS-$$ARCH; \
		if [ $$OS = "windows" ]; then OUTPUT=$$OUTPUT.exe; fi; \
		echo "Building for $$OS/$$ARCH..."; \
		GOOS=$$OS GOARCH=$$ARCH $(GOBUILD) $(LDFLAGS) $(GCFLAGS) $(ASMFLAGS) $(BUILD_TAGS) -o $$OUTPUT .; \
		if [ $$? -eq 0 ]; then \
			echo "✅ Built $$OUTPUT"; \
			ls -lh $$OUTPUT; \
		else \
			echo "❌ Failed to build for $$OS/$$ARCH"; \
		fi; \
	done
	@echo "✅ All builds complete"

## build-linux: Build for Linux (amd64 and arm64)
build-linux: deps
	@echo "🐧 Building for Linux..."
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) $(GCFLAGS) $(ASMFLAGS) $(BUILD_TAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) $(GCFLAGS) $(ASMFLAGS) $(BUILD_TAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 .
	@echo "✅ Linux builds complete"

## build-mac: Build for macOS (Intel and Apple Silicon)
build-mac: deps
	@echo "🍎 Building for macOS..."
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) $(GCFLAGS) $(ASMFLAGS) $(BUILD_TAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) $(GCFLAGS) $(ASMFLAGS) $(BUILD_TAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@echo "✅ macOS builds complete"

## build-windows: Build for Windows
build-windows: deps
	@echo "🪟 Building for Windows..."
	@mkdir -p $(DIST_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) $(GCFLAGS) $(ASMFLAGS) $(BUILD_TAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "✅ Windows build complete"

## install: Install the binary to $GOPATH/bin
install: build
	@echo "📦 Installing $(BINARY_NAME)..."
	$(GOCMD) install $(LDFLAGS) $(GCFLAGS) $(ASMFLAGS) $(BUILD_TAGS) .
	@echo "✅ Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

## uninstall: Remove the binary from $GOPATH/bin
uninstall:
	@echo "🗑️  Uninstalling $(BINARY_NAME)..."
	@rm -f $(shell go env GOPATH)/bin/$(BINARY_NAME)
	@echo "✅ Uninstalled"

## run: Build and run the application
run: build-fast
	@echo "🚀 Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

## clean: Remove build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@echo "✅ Clean complete"

## test: Run tests
test:
	@echo "🧪 Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "✅ Tests complete"

## coverage: Run tests with coverage report
coverage: test
	@echo "📊 Generating coverage report..."
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

## lint: Run linters (requires golangci-lint)
lint:
	@echo "🔍 Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "⚠️  golangci-lint not installed. Run: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin"; \
	fi

## fmt: Format code
fmt:
	@echo "🎨 Formatting code..."
	$(GOFMT) ./...
	@echo "✅ Formatting complete"

## vet: Run go vet
vet:
	@echo "🔎 Running go vet..."
	$(GOVET) ./...
	@echo "✅ Vet complete"

## deps: Download and verify dependencies
deps:
	@echo "📥 Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) verify
	@echo "✅ Dependencies ready"

## tidy: Tidy go.mod
tidy:
	@echo "🧹 Tidying dependencies..."
	$(GOMOD) tidy
	@echo "✅ Dependencies tidied"

## upgrade: Upgrade all dependencies
upgrade:
	@echo "⬆️  Upgrading dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy
	@echo "✅ Dependencies upgraded"

## release: Create a release build (compressed binaries)
release: clean build-all
	@echo "📦 Creating release archives..."
	@mkdir -p $(DIST_DIR)/release
	@cd $(DIST_DIR) && \
	for binary in $(BINARY_NAME)-*; do \
		if [ -f "$binary" ]; then \
			echo "Compressing $binary..."; \
			tar -czf release/$binary.tar.gz $binary; \
			echo "✅ Created release/$binary.tar.gz"; \
		fi; \
	done
	@echo "✅ Release build complete"
	@ls -lh $(DIST_DIR)/release/

## docker: Build Docker image
docker:
	@echo "🐳 Building Docker image..."
	docker build -t $(APP_NAME):$(VERSION) -t $(APP_NAME):latest .
	@echo "✅ Docker image built: $(APP_NAME):$(VERSION)"

## benchmark: Run benchmarks
benchmark:
	@echo "⚡ Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...
	@echo "✅ Benchmarks complete"

## size: Show binary size breakdown
size: build
	@echo "📊 Binary size analysis:"
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME)
	@if command -v du >/dev/null 2>&1; then \
		echo "Disk usage: $(du -h $(BUILD_DIR)/$(BINARY_NAME) | cut -f1)"; \
	fi

## check: Run all checks (fmt, vet, lint, test)
check: fmt vet lint test
	@echo "✅ All checks passed"

## dev: Development mode (watch and rebuild on changes) - requires entr
dev:
	@echo "👀 Watching for changes..."
	@if command -v entr >/dev/null 2>&1; then \
		find . -name '*.go' | entr -r make run; \
	else \
		echo "⚠️  'entr' not installed. Install with: brew install entr (macOS) or apt install entr (Linux)"; \
	fi