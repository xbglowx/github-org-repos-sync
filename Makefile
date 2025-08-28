# Makefile for github-org-repos-sync

# Variables
BINARY_NAME := github-org-repos-sync
MODULE_NAME := github.com/xbglowx/github-org-repos-sync
VERSION_PACKAGE := $(MODULE_NAME)/cmd
BUILD_DIR := ./bin

# Version detection - use git tags if available, otherwise use short git SHA for source builds
GIT_SHA := $(shell git rev-parse --short HEAD 2>/dev/null)
GIT_TAG := $(shell git describe --tags --exact-match 2>/dev/null)
# Dirty detection: Repository is "clean" only if ALL conditions are true:
# 1. No unstaged changes (git diff --quiet returns 0)
# 2. No staged changes (git diff --quiet --cached returns 0)  
# 3. No untracked files (git status --porcelain is empty)
# If ANY condition fails, the && chain fails and executes || echo "-dirty"
GIT_DIRTY := $(shell git diff --quiet 2>/dev/null && git diff --quiet --cached 2>/dev/null && [ -z "$$(git status --porcelain 2>/dev/null)" ] || echo "-dirty")

# If we have an exact tag match, use it; otherwise use git SHA for source builds
ifdef GIT_TAG
    VERSION := $(GIT_TAG)$(GIT_DIRTY)
else ifdef GIT_SHA
    VERSION := $(GIT_SHA)$(GIT_DIRTY)
else
    VERSION := dev-unknown
endif

# Go build flags
LDFLAGS := -X $(VERSION_PACKAGE).Version=$(VERSION)
BUILD_FLAGS := -ldflags "$(LDFLAGS)"

# Platform detection for local builds
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
GO_VERSION := $(shell go version)

# Cross-compilation targets
PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/arm64

.PHONY: help build build-all test version info clean

# Default target
all: test build

# Help target
help:
	@printf "%s\n" \
		"Available targets:" \
		"  build      - Build binary for current platform" \
		"  build-all  - Build binaries for all platforms" \
		"  test       - Run tests" \
		"  clean      - Clean build artifacts" \
		"  version    - Show version that would be built" \
		"  info       - Show build information" \
		"  help       - Show this help message"

# Show version
version:
	@echo "Version: $(VERSION)"

# Build for current platform
build:
	@echo "Building $(BINARY_NAME) v$(VERSION) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

# Build for all platforms
build-all:
	@echo "Building $(BINARY_NAME) v$(VERSION) for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@$(foreach platform,$(PLATFORMS), \
		echo "Building for $(platform)..."; \
		GOOS=$(word 1,$(subst /, ,$(platform))) \
		GOARCH=$(word 2,$(subst /, ,$(platform))) \
		go build $(BUILD_FLAGS) \
			-o $(BUILD_DIR)/$(BINARY_NAME)-$(word 1,$(subst /, ,$(platform)))-$(word 2,$(subst /, ,$(platform)))$(if $(findstring windows,$(platform)),.exe) . && \
	) true

# Run tests
test:
	@echo "Running tests..."
	go test -v -race ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

# Print build information
info:
	@printf "%s\n" \
		"Build Information:" \
		"  Binary Name: $(BINARY_NAME)" \
		"  Module:      $(MODULE_NAME)" \
		"  Version:     $(VERSION)" \
		"  Go Version:  $(GO_VERSION)" \
		"  Platform:    $(GOOS)/$(GOARCH)" \
		"  Build Dir:   $(BUILD_DIR)"
