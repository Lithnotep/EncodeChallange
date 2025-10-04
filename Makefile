# Makefile for Encode Challenge Go Project

.PHONY: build run test clean help

# Build the application
build:
	go build -o bin/encodechallenge .

# Run the application
run:
	go run .

# Run tests
test:
	go test -v .

# Run tests with coverage
test-coverage:
	go test -v -cover .

# Clean build artifacts
clean:
	rm -rf bin/

# Format code
fmt:
	go fmt .

# Run go mod tidy
tidy:
	go mod tidy

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Clean build artifacts"
	@echo "  fmt           - Format code"
	@echo "  tidy          - Run go mod tidy"
	@echo "  help          - Show this help message"