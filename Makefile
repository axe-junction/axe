# Axe Server - Microservices Management
.PHONY: help build up down logs clean dev prod test lint

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development commands
dev: ## Start development environment with hot reloading
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build

dev-detached: ## Start development environment in background
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build

# Production commands
prod: ## Start production environment
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build

prod-with-monitoring: ## Start production with monitoring stack
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml --profile monitoring up -d --build

# Basic commands
build: ## Build all services
	docker-compose build

up: ## Start all services (basic setup)
	docker-compose up -d

down: ## Stop all services
	docker-compose down

down-volumes: ## Stop all services and remove volumes
	docker-compose down -v

# Monitoring and debugging
logs: ## Show logs for all services
	docker-compose logs -f

logs-routing: ## Show logs for routing service
	docker-compose logs -f routing-service

logs-vtc: ## Show logs for VTC service
	docker-compose logs -f vtc-service

logs-db: ## Show logs for database
	docker-compose logs -f postgres

# Individual services
routing-only: ## Start only routing service with its dependencies
	docker-compose up -d routing-db routing-service

vtc-only: ## Start only VTC service with its dependencies
	docker-compose up -d vtc-db vtc-service

data-only: ## Start only data API service with its dependencies
	docker-compose up -d data-db data-api

# Database management
init-dbs: ## Initialize all databases with schema and sample data
	@echo "Waiting for databases to be ready..."
	@sleep 10
	docker-compose exec routing-db psql -U routing_user -d routing_db -f /docker-entrypoint-initdb.d/init-routing-db.sql
	docker-compose exec vtc-db psql -U vtc_user -d vtc_db -f /docker-entrypoint-initdb.d/init-vtc-db.sql
	docker-compose exec data-db psql -U data_user -d data_db -f /docker-entrypoint-initdb.d/init-data-db.sql

reset-dbs: ## Reset all databases (drops and recreates them)
	docker-compose stop routing-db vtc-db data-db
	docker-compose rm -f routing-db vtc-db data-db
	docker volume rm axe-server_routing_db_data axe-server_vtc_db_data axe-server_data_db_data 2>/dev/null || true
	docker-compose up -d routing-db vtc-db data-db

backup-dbs: ## Backup all databases
	mkdir -p backups
	docker-compose exec routing-db pg_dump -U routing_user routing_db > backups/routing_db_$(shell date +%Y%m%d_%H%M%S).sql
	docker-compose exec vtc-db pg_dump -U vtc_user vtc_db > backups/vtc_db_$(shell date +%Y%m%d_%H%M%S).sql
	docker-compose exec data-db pg_dump -U data_user data_db > backups/data_db_$(shell date +%Y%m%d_%H%M%S).sql

# Database-specific logs
logs-routing-db: ## Show logs for routing database
	docker-compose logs -f routing-db

logs-vtc-db: ## Show logs for VTC database
	docker-compose logs -f vtc-db

logs-data-db: ## Show logs for data database
	docker-compose logs -f data-db

# Database connections (useful for development)
connect-routing-db: ## Connect to routing database
	docker-compose exec routing-db psql -U routing_user -d routing_db

connect-vtc-db: ## Connect to VTC database
	docker-compose exec vtc-db psql -U vtc_user -d vtc_db

connect-data-db: ## Connect to data database
	docker-compose exec data-db psql -U data_user -d data_db

# PGAdmin access
pgadmin: ## Open PGAdmin web interface
	@echo "PGAdmin is available at: http://localhost:5050"
	@echo "Email: admin@axe.com"
	@echo "Password: admin"

# Testing
test: ## Run tests for all services
	docker-compose -f docker-compose.yml -f docker-compose.test.yml up --build --abort-on-container-exit

test-routing: ## Run tests for routing service only
	cd packages/routing-service && go test ./...

test-vtc: ## Run tests for VTC service only
	cd packages/vtc-service && go test ./...

# Code quality
lint: ## Run linters for all Go services
	cd packages/routing-service && golangci-lint run
	cd packages/vtc-service && golangci-lint run

fmt: ## Format code for all Go services
	cd packages/routing-service && go fmt ./...
	cd packages/vtc-service && go fmt ./...

# Cleanup
clean: ## Clean up Docker resources
	docker-compose down -v --remove-orphans
	docker system prune -f
	docker volume prune -f

clean-all: ## Nuclear cleanup - remove everything
	docker-compose down -v --remove-orphans
	docker system prune -a -f
	docker volume prune -f

# Health checks
health: ## Check health of all services
	@echo "Checking service health..."
	@curl -f http://localhost:3001/health && echo "✅ Routing service healthy" || echo "❌ Routing service down"
	@curl -f http://localhost:8080/health && echo "✅ VTC service healthy" || echo "❌ VTC service down"
	@curl -f http://localhost:8000/health && echo "✅ Data API healthy" || echo "❌ Data API down"

# Monitoring
monitor: ## Start monitoring stack only
	docker-compose --profile monitoring up -d prometheus grafana

# Backup
backup-db: ## Backup database
	docker-compose exec postgres pg_dump -U axe_user axe_db > backup_$(shell date +%Y%m%d_%H%M%S).sql

restore-db: ## Restore database from backup (usage: make restore-db FILE=backup.sql)
	docker-compose exec -T postgres psql -U axe_user -d axe_db < $(FILE)

# Environment setup
setup-dev: ## Setup development environment
	@echo "Setting up development environment..."
	@cp .env.example .env 2>/dev/null || echo "No .env.example found"
	@echo "Installing development dependencies..."
	cd packages/routing-service && go mod download
	cd packages/vtc-service && go mod download
	@echo "✅ Development environment ready!"

# API testing
test-api: ## Test API endpoints
	@echo "Testing routing service..."
	curl -s "http://localhost:3001/routing/best?fromLat=36.75&fromLng=3.05&toLat=36.76&toLng=3.06" | jq .
	@echo "Testing VTC service..."
	curl -s "http://localhost:8080/api/estimate" | jq . || echo "VTC service endpoint may need parameters"

# Documentation
docs: ## Generate API documentation
	@echo "API documentation available at:"
	@echo "- Routing Service: http://localhost:3001/"
	@echo "- VTC Service: http://localhost:8080/"
	@echo "- Data API: http://localhost:8000/docs"
	@echo "- Traefik Dashboard: http://localhost:8090/ (if proxy profile enabled)"
	@echo "- Grafana: http://localhost:3000/ (if monitoring profile enabled)"

# Version info
version: ## Show version information
	@echo "Axe Server Microservices"
	@echo "========================"
	@docker-compose version
	@echo ""
	@echo "Services:"
	@docker-compose ps
