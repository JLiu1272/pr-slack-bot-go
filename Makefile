# Specify BINARY_NAME variable
BINARY_NAME=bin/server

# Create a command to build the binary for a go project
# Usage: make build
build:
	@echo "Building binary..."
	@go build -o $(BINARY_NAME) -v

# Run playground. A fast way to experiment 
# Usage: make mini-playground
mini-playground:
	@echo "Running playground..."
	@go run ./playground/main.go

# Usage: make start 
start:
	@echo "Starting server..."
	@go run ./actions.go ./blockBuilder.go ./helpMessage.go ./listPRs.go ./server.go ./usernameExist.go 