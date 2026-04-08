.PHONY: all build run dev test clean fmt vet help

APP_NAME=student-api
BINARY_DIR=bin
BINARY_PATH=$(BINARY_DIR)/$(APP_NAME)
MAIN_FILE=cmd/server/main.go

# Default target: running 'make' will run 'make build'
all: build

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build       Build the application binary"
	@echo "  run         Run the compiled binary"
	@echo "  dev         Run the app directly with 'go run' (best for development)"
	@echo "  test        Run all unit tests"
	@echo "  fmt         Format the source code using 'go fmt'"
	@echo "  vet         Analyze the source code for common errors using 'go vet'"
	@echo "  clean       Remove build artifacts (bin/ directory)"

build:
	@echo "Building application..."
	@mkdir -p $(BINARY_DIR)
	go build -o $(BINARY_PATH) $(MAIN_FILE)

run:
	@echo "Running compiled binary..."
	@if [ -f $(BINARY_PATH) ]; then \
		./$(BINARY_PATH); \
	else \
		echo "Binary not found. Run 'make build' first."; \
	fi

dev:
	@echo "Running application with 'go run'..."
	go run $(MAIN_FILE)

test:
	@echo "Running tests..."
	go test -v ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...

vet:
	@echo "Vetting code..."
	go vet ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BINARY_DIR)
