.PHONY: build run test clean fmt help

BINARY_NAME=crimson
BUILD_DIR=bin

build:
	@echo "🔨 Building Crimson..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/crimson/main.go
	@echo "✅ Build complete"

run:
	@go run cmd/crimson/main.go

test:
	@go test -v ./...

fmt:
	@go fmt ./...

clean:
	@go clean
	@rm -rf $(BUILD_DIR)

help:
	@echo "build   → Build binary"
	@echo "run     → Run server"
	@echo "test    → Run tests"
	@echo "fmt     → Format code"
	@echo "clean   → Clean artifacts"