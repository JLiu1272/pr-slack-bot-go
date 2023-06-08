# Specify BINARY_NAME variable
BINARY_NAME=bin/server

# Create a command to build the binary for a go project
# Usage: make build
build:
	@echo "Building binary..."
	@go build -o $(BINARY_NAME) -v

start:
	@echo "Starting server..."
	@go run *.go 