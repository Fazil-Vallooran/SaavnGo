.PHONY: build run test clean dev

# Build the application
build:
	@echo "Building application..."
	go build -o bin/jiosaavn-api

# Run the application
run:
	@echo "Running application..."
	go run main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f jiosaavn-api

# Run in development mode with hot reload (requires air)
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Install it with: go install github.com/cosmtrek/air@latest"; \
		echo "Running without hot reload..."; \
		go run main.go; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run the application on a different port
run-alt:
	@echo "Running application on port 3000..."
	SERVER_PORT=3000 go run main.go

# Build for multiple platforms
build-all:
	@echo "Building for all platforms..."
	GOOS=linux GOARCH=amd64 go build -o bin/jiosaavn-api-linux-amd64
	GOOS=windows GOARCH=amd64 go build -o bin/jiosaavn-api-windows-amd64.exe
	GOOS=darwin GOARCH=amd64 go build -o bin/jiosaavn-api-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o bin/jiosaavn-api-darwin-arm64

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  run-alt       - Run the application on port 3000"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  dev           - Run in development mode with hot reload"
	@echo "  fmt           - Format code"
	@echo "  lint          - Run linter"
	@echo "  deps          - Install dependencies"
	@echo "  build-all     - Build for all platforms"
	@echo "  help          - Show this help message"
