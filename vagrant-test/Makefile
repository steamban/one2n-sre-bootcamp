.PHONY: deploy deploy-stop deploy-restart deploy-logs ps help

# Load .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

COMPOSE_FILE = deployment/docker-compose.yml

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  deploy          Start all containers in detached mode"
	@echo "  deploy-stop     Stop all containers"
	@echo "  deploy-restart  Restart all containers"
	@echo "  deploy-logs     View logs from all containers"
	@echo "  ps              Show running containers"

deploy:
	docker-compose -f $(COMPOSE_FILE) up -d
	@echo "Services started. API available at http://localhost:8080"

deploy-stop:
	docker-compose -f $(COMPOSE_FILE) down

deploy-restart:
	docker-compose -f $(COMPOSE_FILE) restart

deploy-logs:
	docker-compose -f $(COMPOSE_FILE) logs -f

ps:
	docker-compose -f $(COMPOSE_FILE) ps