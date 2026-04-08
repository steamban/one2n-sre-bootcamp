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
Copy the example environment file and adjust as needed:
```bash
cp .env.example .env 
```

Ensure your `.env` file contains:
```env
PORT=8080
DATABASE_PATH=students.db
```

## Running the Application

### Using Go Directly
```bash
go run cmd/server/main.go
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
