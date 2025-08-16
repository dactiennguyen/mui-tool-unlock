# MUI Tool Unlock - Makefile
# Build both GUI and Terminal versions

# Variables
BINARY_NAME_GUI=mui-tool-unlock
BINARY_NAME_TERMINAL=mui-tool-unlock-terminal
BUILD_DIR=build
DIST_DIR=dist

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w"
CGO_ENABLED=1

# Version info
VERSION=1.0.0
BUILD_TIME=$(shell date +%Y-%m-%d_%H:%M:%S)
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Targets
.PHONY: all build-gui build-terminal clean install test cross-compile help

# Default target - build both versions
all: clean build-gui build-terminal

# Build GUI version
build-gui:
	@echo "🖥️  Building GUI version..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) $(GOBUILD) -tags gui $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME_GUI) .
	@echo "✅ GUI build completed: $(BUILD_DIR)/$(BINARY_NAME_GUI)"

# Build Terminal version  
build-terminal:
	@echo "⌨️  Building Terminal version..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 $(GOBUILD) -tags terminal $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME_TERMINAL) .
	@echo "✅ Terminal build completed: $(BUILD_DIR)/$(BINARY_NAME_TERMINAL)"

# Cross-compilation targets
cross-compile: clean
	@echo "🌍 Cross-compiling for multiple platforms..."
	@mkdir -p $(DIST_DIR)
	
	# Windows GUI (requires CGO for Fyne)
	@echo "🪟 Building Windows GUI..."
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 $(GOBUILD) -tags gui $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME_GUI)-windows-amd64.exe .
	
	# Windows Terminal
	@echo "🪟 Building Windows Terminal..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -tags terminal $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME_TERMINAL)-windows-amd64.exe .
	
	# macOS GUI (requires CGO for Fyne)
	@echo "🍎 Building macOS GUI..."
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(GOBUILD) -tags gui $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME_GUI)-darwin-amd64 .
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 $(GOBUILD) -tags gui $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME_GUI)-darwin-arm64 .
	
	# macOS Terminal
	@echo "🍎 Building macOS Terminal..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -tags terminal $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME_TERMINAL)-darwin-amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) -tags terminal $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME_TERMINAL)-darwin-arm64 .
	
	# Linux GUI (requires CGO for Fyne)
	@echo "🐧 Building Linux GUI..."
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -tags gui $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME_GUI)-linux-amd64 .
	
	# Linux Terminal
	@echo "🐧 Building Linux Terminal..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -tags terminal $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME_TERMINAL)-linux-amd64 .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -tags terminal $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME_TERMINAL)-linux-arm64 .
	
	@echo "✅ Cross-compilation completed! Check $(DIST_DIR)/ directory"

# Quick build for current platform
quick: build-gui build-terminal
	@echo "🚀 Quick build completed!"

# Build and run GUI version
run-gui: build-gui
	@echo "🚀 Running GUI version..."
	./$(BUILD_DIR)/$(BINARY_NAME_GUI)

# Build and run Terminal version
run-terminal: build-terminal
	@echo "🚀 Running Terminal version..."
	./$(BUILD_DIR)/$(BINARY_NAME_TERMINAL)

# Test Terminal with parameters
test-terminal: build-terminal
	@echo "🧪 Testing Terminal version..."
	./$(BUILD_DIR)/$(BINARY_NAME_TERMINAL) --help
	@echo ""
	./$(BUILD_DIR)/$(BINARY_NAME_TERMINAL) --version

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)
	rm -f $(BINARY_NAME_GUI) $(BINARY_NAME_TERMINAL)
	@echo "✅ Clean completed!"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✅ Dependencies installed!"

# Run tests
test:
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...

# Install binaries to system (requires sudo on Unix)
install: build-gui build-terminal
	@echo "📥 Installing binaries..."
	@if [ "$(shell uname)" = "Darwin" ] || [ "$(shell uname)" = "Linux" ]; then \
		echo "Installing to /usr/local/bin/ (requires sudo)..."; \
		sudo cp $(BUILD_DIR)/$(BINARY_NAME_GUI) /usr/local/bin/; \
		sudo cp $(BUILD_DIR)/$(BINARY_NAME_TERMINAL) /usr/local/bin/; \
		sudo chmod +x /usr/local/bin/$(BINARY_NAME_GUI); \
		sudo chmod +x /usr/local/bin/$(BINARY_NAME_TERMINAL); \
	else \
		echo "Manual installation required on Windows"; \
		echo "Copy $(BUILD_DIR)/$(BINARY_NAME_GUI).exe and $(BUILD_DIR)/$(BINARY_NAME_TERMINAL).exe to your PATH"; \
	fi
	@echo "✅ Installation completed!"

# Show file sizes
sizes: build-gui build-terminal
	@echo "📊 Binary sizes:"
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME_GUI) $(BUILD_DIR)/$(BINARY_NAME_TERMINAL)

# Package for distribution
package: cross-compile
	@echo "📦 Creating distribution packages..."
	@mkdir -p $(DIST_DIR)/packages
	cd $(DIST_DIR) && tar -czf packages/$(BINARY_NAME_GUI)-$(VERSION)-windows-amd64.tar.gz $(BINARY_NAME_GUI)-windows-amd64.exe $(BINARY_NAME_TERMINAL)-windows-amd64.exe
	cd $(DIST_DIR) && tar -czf packages/$(BINARY_NAME_GUI)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME_GUI)-darwin-amd64 $(BINARY_NAME_TERMINAL)-darwin-amd64
	cd $(DIST_DIR) && tar -czf packages/$(BINARY_NAME_GUI)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME_GUI)-darwin-arm64 $(BINARY_NAME_TERMINAL)-darwin-arm64
	cd $(DIST_DIR) && tar -czf packages/$(BINARY_NAME_GUI)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME_GUI)-linux-amd64 $(BINARY_NAME_TERMINAL)-linux-amd64
	cd $(DIST_DIR) && tar -czf packages/$(BINARY_NAME_GUI)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME_GUI)-linux-arm64 $(BINARY_NAME_TERMINAL)-linux-arm64
	@echo "✅ Packages created in $(DIST_DIR)/packages/"

# Development mode - rebuild on file changes (requires entr)
dev-gui:
	@echo "🔄 Development mode for GUI (auto-rebuild)..."
	@if command -v entr >/dev/null 2>&1; then \
		find . -name "*.go" | entr -r make run-gui; \
	else \
		echo "Install 'entr' for auto-rebuild: brew install entr"; \
		make run-gui; \
	fi

dev-terminal:
	@echo "🔄 Development mode for Terminal (auto-rebuild)..."
	@if command -v entr >/dev/null 2>&1; then \
		echo "./terminal.go" | entr -r make run-terminal; \
	else \
		echo "Install 'entr' for auto-rebuild: brew install entr"; \
		make run-terminal; \
	fi

# Show build info
info:
	@echo "ℹ️  Build Information:"
	@echo "   Version: $(VERSION)"
	@echo "   Build Time: $(BUILD_TIME)"
	@echo "   Git Commit: $(GIT_COMMIT)"
	@echo "   Go Version: $(shell $(GOCMD) version)"

# Help target
help:
	@echo "🔧 MUI Tool Unlock - Makefile Help"
	@echo "=================================="
	@echo ""
	@echo "Main Targets:"
	@echo "  all                Build both GUI and Terminal versions"
	@echo "  build-gui          Build GUI version only"
	@echo "  build-terminal     Build Terminal version only"
	@echo "  cross-compile      Build for all platforms"
	@echo "  clean              Clean build artifacts"
	@echo ""
	@echo "Run Targets:"
	@echo "  run-gui            Build and run GUI version"
	@echo "  run-terminal       Build and run Terminal version"  
	@echo "  test-terminal      Build and test Terminal version"
	@echo ""
	@echo "Development:"
	@echo "  dev-gui            Development mode with auto-rebuild (GUI)"
	@echo "  dev-terminal       Development mode with auto-rebuild (Terminal)"
	@echo "  deps               Install dependencies"
	@echo "  test               Run tests"
	@echo ""
	@echo "Distribution:"
	@echo "  package            Create distribution packages"
	@echo "  install            Install binaries to system"
	@echo "  sizes              Show binary sizes"
	@echo ""
	@echo "Utilities:"
	@echo "  info               Show build information"
	@echo "  help               Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make                    # Build both versions"
	@echo "  make build-gui          # Build GUI only"
	@echo "  make cross-compile      # Build for all platforms"
	@echo "  make run-terminal       # Run Terminal version"
