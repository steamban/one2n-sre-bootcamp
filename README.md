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

## Vagrant Setup (Local Development with UTM)

This project includes a Vagrant configuration for running the full stack locally using UTM on Mac.

### Prerequisites
- **UTM** (or another VM hypervisor that supports Ubuntu)
- **~4GB RAM** allocated to the VM
- **2-4 CPU cores** allocated to the VM

### Installing Required Software

If you don't have these installed, run the following commands:

### Install UTM
```bash
brew install --cask utm
```

### Install Vagrant
```bash
brew tap hashicorp/tap
brew install hashicorp/tap/hashicorp-vagrant
```

### Install Vagrant UTM Plugin
```bash
vagrant plugin install vagrant_utm
```

### Getting Started

1. **Navigate to the vagrant directory:**
   ```bash
   cd vagrant-test
   ```

2. **Start the VM:**
   ```bash
   vagrant up
   ```

3. **SSH into the VM:**
   ```bash
   vagrant ssh
   ```

4. **Start the application stack:**
   ```bash
   cd /vagrant
   make deploy
   ```

5. **Access the application:**
   - **API (via nginx load balancer)**: `http://localhost:8080`
   - **Healthcheck**: `http://localhost:8080/healthcheck`


### Stopping the Stack

```bash
cd deployment
make deploy-stop
```

To stop and remove the VM:

```bash
vagrant halt        # Stop the VM (preserve data)
vagrant destroy      # Destroy the VM (remove all data)
```

### Troubleshooting

- **VM not starting**: Ensure UTM/vagrant is properly installed and virtualization is enabled
- **Port conflicts**: Ensure nothing else is running on port 8080
- **Database connection issues**: Check that postgres container is healthy with `docker compose ps`

## Helm Setup (Kubernetes Deployment)

This project includes a Helm chart for deploying the Student API to Kubernetes along with its infrastructure dependencies (PostgreSQL, Vault, and External Secrets Operator).

### Prerequisites
- **Kubernetes cluster** (e.g., kind, minikube, or a cloud-based cluster)
- **Helm 3+**
- **kubectl** configured to access your cluster
- **helmfile** (optional, for declarative deployments)

For setting up a local minikube cluster, you can use the provided `infra.sh` script:
```bash
cd helm-setup
./infra.sh
```

### Infrastructure Components

The helmfile deploys three main components:
1. **Vault** - For secrets management (runs in `vault` namespace)
2. **External Secrets Operator (ESO)** - For syncing secrets from Vault to Kubernetes (runs in `external-secrets` namespace)
3. **Student API** - The application with PostgreSQL as a dependency (runs in `student-api` namespace)

### Deploying with Helmfile

1. **Navigate to the helm setup directory:**
   ```bash
   cd helm-setup
   ```

2. **Deploy infra components:**
   ```bash
   helmfile apply --selector tier=infra
   ```

3. **Deploy app component:**
   ```bash
   helmfile apply --selector tier=app
   ```

4. **Verify the deployment:**
   ```bash
   kubectl get pods -n student-api
   kubectl get pods -n vault
   kubectl get pods -n external-secrets
   ```

### Accessing the Application

Once deployed, port-forward to access locally:
```bash
kubectl port-forward svc/student-api-student-api 8080:8080 -n student-api
```

### Verifying the Installation

Check the healthcheck endpoint:
```bash
curl http://localhost:8080/healthcheck
```

### Cleanup

To uninstall all resources:
```bash
helmfile destroy
```


## Project Structure
- `cmd/server/`: Entry point for the application.
- `internal/`: Private application and library code.
- `pkg/`: Public library code.
- `migrations/`: SQL migration files.
- `helm-setup/`: Run application with helm
