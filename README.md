# Student CRUD REST API

## Purpose
A simple REST API webserver for student management, built with Golang and the Gin framework. This project follows the Twelve-Factor App methodology and uses SQLite for data persistence.

## Prerequisites
- **Go 1.22+**:
- **SQLite**: 

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
cp .env.example .env 
```

Edit the `.env` file to set necessary variables. Key variables include:
- `PORT`: The port on which the API server will listen (e.g., `8080`).
- `DATABASE_PATH`: The path to the SQLite database file (defaults to `students.db`).

## Database Migrations
Database schema changes are managed using `golang-migrate`. The migration files are located in the `migrations/` directory. You can apply or rollback migrations using the `make` commands:
- `make migrate-up`: Applies all pending database migrations.
- `make migrate-down`: Rolls back the last applied database migration.
- `make migrate-new name=<name>`: Creates a new pair of migration files (up and down) for a new schema change. Replace `<name>` with a descriptive name for your migration.`

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
