.PHONY: build test test-race test-cover test-cover-html lint fmt vet install clean ci run

# Build
build:
	go build -o devbot ./cmd/devbot

install:
	go install ./cmd/devbot

# Testing
test:
	go test ./... -v

test-race:
	go test ./... -v -race

test-cover:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -func=coverage.out | tail -1

test-cover-html:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Quality
fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	@which golangci-lint > /dev/null || (echo "Install: brew install golangci-lint" && exit 1)
	golangci-lint run

# CI (runs everything)
ci: fmt vet test-race test-cover
	@echo "CI passed"

# Development
run: build
	./devbot status

# Cleanup
clean:
	rm -f devbot coverage.out coverage.html
