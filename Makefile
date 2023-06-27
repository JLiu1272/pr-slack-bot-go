# Specify BINARY_NAME variable
BINARY_NAME=bin/server

# Create a command to build the binary for a go project
# Usage: make build
build:
	@echo "Building binary..."
	@go build -o $(BINARY_NAME) -v

# Using a crude way to start a server. This is not a recommended approach 
# Usage: make start 
start:
	@echo "Starting server..."
	@go run ./*.go 