# Student CRUD REST API

## Purpose
A simple REST API webserver for student management, built with Golang and the Gin framework. This project uses PostgreSQL for data persistence.

## Prerequisites
- **Go 1.22+**
- **PostgreSQL 14+**
- **Docker** (for containerized setup)

### 1. Clone the repository
```bash
git clone <repository-url>
cd student-api
```

### 2. Install dependencies
```bash
go mod tidy
```

### 3. Configure Environment Variables
The application uses environment variables for configuration. Copy the example file to create your local `.env`:
```bash
cp example.env .env 
```

Edit the `.env` file to set necessary variables. Key variables include:
- `PORT`: The port on which the API server will listen (e.g., `8080`).
- `DB_HOST`: PostgreSQL database host (e.g., `localhost`).
- `DB_PORT`: PostgreSQL database port (e.g., `5432`).
- `DB_USER`: PostgreSQL database user.
- `DB_PASSWORD`: PostgreSQL database password.
- `DB_NAME`: PostgreSQL database name.
- `DB_SSLMODE`: PostgreSQL SSL mode (e.g., `disable`).

## Database Migrations
Database schema changes are managed using `golang-migrate`. The migration files are located in the `migrations/` directory. You can apply or rollback migrations using the `make` commands:
- `make migrate-up`: Applies all pending database migrations.
- `make migrate-down`: Rolls back the last applied database migration.
- `make migrate-new name=<name>`: Creates a new pair of migration files (up and down) for a new schema change. Replace `<name>` with a descriptive name for your migration.

## Running the Application

### Using Go Directly
```bash
go run cmd/server/main.go
```

### Using Make Commands
This project includes a `Makefile` to streamline common development tasks. Below are the available commands:

- `make build`: Builds all application binaries (api server and migration tool).
- `make run`: Builds and runs the API server.
- `make dev`: Runs the API server directly using `go run` (useful for quick development cycles).
- `make test`: Executes all Go tests in the project.
- `make fmt`: Formats the Go source code.
- `make vet`: Analyzes the Go source code for potential errors.
- `make clean`: Removes all built binaries and artifacts.
- `make migrate-new name=<name>`: Creates a new pair of migration files (up and down) with the given name.
- `make migrate-up`: Applies all pending database migrations.
- `make migrate-down`: Rolls back the last applied database migration.

### Using Docker
Start the database and run the application in containers:
```bash
make docker-build      # Create the app container image
make docker-pg-start   # Start PostgreSQL container
make docker-migrate    # Run database migrations
make docker-run        # Start the API container
```

Stop the containers:
```bash
make docker-app-stop   # Stop API container
make docker-pg-stop     # Stop and remove PostgreSQL container
```

### Using Docker Compose
Docker Compose commands are also available as Make targets for convenience:
```bash
make docker-compose-up      # Start all services (builds if needed)
make docker-compose-down    # Stop and remove containers
make docker-compose-clean   # Stop and remove containers with volumes
```

## Verifying the Installation
Once the server is running, you can test the healthcheck endpoint:
```bash
curl -i http://localhost:8080/healthcheck
```
You should receive a `200 OK` response with `{"status":"UP"}`.

## Project Structure
- `cmd/server/`: Entry point for the application.
- `internal/`: Private application and library code.
- `pkg/`: Public library code.
- `migrations/`: SQL migration files.
