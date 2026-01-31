.PHONY: build test lint dev clean setup

# Build the application
build:
	go build -o app ./cmd/app

# Run all tests
test:
	go test ./... -v

# Run linter
lint:
	golangci-lint run

# Start development server
dev:
	go run ./cmd/app

# Clean build artifacts
clean:
	rm -f app

# First time setup
setup:
	go mod download
