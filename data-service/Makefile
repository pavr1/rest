# Bar-Restaurant Data Service Makefile
# Docker-related commands only

.PHONY: start stop restart logs connect status fresh clean migrate migrate-down migrate-status test help

.DEFAULT_GOAL := help

# =============================================================================
# 🚀 MAIN COMMANDS
# =============================================================================

start: _network ## Start containers (PostgreSQL + PgAdmin)
	@echo "🍺 Starting Bar-Restaurant Data Service..."
	@cd docker && docker-compose up -d --build
	@echo "⏳ Waiting for database..."
	@sleep 3
	@echo "✅ Started!"
	@make _info

stop: ## Stop all containers
	@echo "🛑 Stopping containers..."
	@cd docker && docker-compose down
	@echo "✅ Stopped!"

restart: stop start ## Restart database container

logs: ## View logs (use: make logs or make logs s=postgres)
	@cd docker && docker-compose logs -f $(s)

connect: ## Connect to PostgreSQL CLI
	@docker exec -it barrest_postgres psql -U postgres -d barrest_db

status: ## Show container status
	@cd docker && docker-compose ps

fresh: clean start ## Clean everything and start fresh
	@echo "⏳ Waiting for database to be ready..."
	@sleep 2
	@make migrate
	@echo "🎉 Fresh install complete!"

clean: ## Remove containers and volumes
	@cd docker && docker-compose down -v --remove-orphans 2>/dev/null || true
	@echo "✅ Data Service cleaned!"

migrate: ## Apply pending migrations
	@./scripts/migrate.sh up

migrate-down: ## Rollback last migration
	@./scripts/migrate.sh down

migrate-status: ## Show applied migrations
	@./scripts/migrate.sh status

test: ## Run unit tests
	@echo "🧪 Running tests..."
	@go test -v ./...

help: ## Show this help
	@echo "🍺 Bar-Restaurant Data Service"
	@echo ""
	@echo "Usage: make [command]"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'
	@echo ""
	@echo "Quick Start: make fresh"

# =============================================================================
# 🔧 INTERNAL
# =============================================================================

_network:
	@docker network create docker_barrest_network 2>/dev/null || true

_info:
	@echo ""
	@echo "📝 Connection Info:"
	@echo "   PostgreSQL: localhost:5432 (postgres/postgres123)"
	@echo "   Database:   barrest_db"
	@echo "   PgAdmin:    http://localhost:8080 (admin@barrest.com / admin123)"