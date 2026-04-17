.PHONY: all build run test clean fmt vet help migrate-new migrate-up migrate-down dev docker-build docker-pg-start docker-pg-stop docker-migrate docker-run docker-app-stop

APP_NAME=student-api
MIGRATE_NAME=student-migrate
BINARY_DIR=bin
SERVER_MAIN=cmd/server/main.go
MIGRATE_MAIN=cmd/migrate/main.go
VERSION?=0.1.0

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
	@echo "  docker-build       Build Docker image (requires VERSION=x.y.z)"
	@echo "  docker-pg-start    Start Postgres container"
	@echo "  docker-pg-stop     Stop and remove Postgres container"
	@echo "  docker-migrate     Run migrations (requires docker-pg-start)"
	@echo "  docker-run         Run API container (requires docker-pg-start)"
	@echo "  docker-app-stop    Stop and remove API container"

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
		last_sim=$$(ls migrations/*.up.sql 2>/dev/null | sort | tail -n 1 | grep -oE '[0-9]+' | head -1 || echo 0); \
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

docker-build:
	docker build -t student-api:$(VERSION) .

docker-pg-start:
	docker run -d --name student-postgres \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=secret \
		-e POSTGRES_DB=students \
		-v student-pgdata:/var/lib/postgresql/data \
		-p 5433:5432 \
		postgres:16-alpine

docker-pg-stop:
	docker stop student-postgres && docker rm student-postgres

docker-migrate: 
	docker run --rm \
		-e DB_HOST=host.docker.internal \
		-e DB_PORT=5433 \
		-e DB_USER=postgres \
		-e DB_PASSWORD=secret \
		-e DB_NAME=students \
		-e DB_SSLMODE=disable \
		student-api:$(VERSION) ./migrate up

docker-run: 
	docker run \
		-e DB_HOST=host.docker.internal \
		-e DB_PORT=5433 \
		-e DB_USER=postgres \
		-e DB_PASSWORD=secret \
		-e DB_NAME=students \
		-e DB_SSLMODE=disable \
		-p 8080:8080 \
		student-api:$(VERSION)

docker-app-stop:
	docker stop student-api && docker rm student-api

docker-compose-up:
	docker compose up

docker-compose-down:
	docker compose down

docker-compose-clean:
	docker compose down -v
	docker rmi $$(docker images "student-api" -q)