.PHONY: all build run test clean fmt vet help migrate-new migrate-up migrate-down dev

APP_NAME=student-api
MIGRATE_NAME=student-migrate
BINARY_DIR=bin
SERVER_MAIN=cmd/server/main.go
MIGRATE_MAIN=cmd/migrate/main.go

# Default target
all: build

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build              Build all binaries"
	@echo "  run                Run the API server"	
	@echo "  dev                Run the API server with go run"
	@echo "  test               Run all tests"
	@echo "  fmt                Format the source code"
	@echo "  vet                Analyze the source code"
	@echo "  clean              Remove build artifacts"
	@echo "  migrate-new name=<n> Create a new migration file"
	@echo "  migrate-up         Apply all pending migrations"
	@echo "  migrate-down       Rollback the last migration"

build:
	@echo "Building binaries..."
	@mkdir -p $(BINARY_DIR)
	go build -o $(BINARY_DIR)/$(APP_NAME) $(SERVER_MAIN)
	go build -o $(BINARY_DIR)/$(MIGRATE_NAME) $(MIGRATE_MAIN)

run: build
	@echo "Starting API server..."
	./$(BINARY_DIR)/$(APP_NAME)

dev:
	@echo "Running API server with 'go run'..."
	go run $(SERVER_MAIN)

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

migrate-new:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-new name=your_migration_name"; \
	else \
		last_sim=$$(ls migrations/*.up.sql 2>/dev/null | sort | tail -n 1 | grep -oE '^[0-9]+' || echo 0); \
		next_num=$$(printf "%06d" $$(($$last_sim + 1))); \
		touch migrations/$${next_num}_$(name).up.sql; \
		touch migrations/$${next_num}_$(name).down.sql; \
		echo "Created migrations/$${next_num}_$(name).up.sql and migrations/$${next_num}_$(name).down.sql"; \
	fi

migrate-up: build
	@echo "Running migrations up..."
	./$(BINARY_DIR)/$(MIGRATE_NAME) up

migrate-down: build
	@echo "Running migrations down..."
	./$(BINARY_DIR)/$(MIGRATE_NAME) down
