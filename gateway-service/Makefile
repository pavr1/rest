# Gateway Service Makefile
# Docker-related commands + test

.PHONY: start stop restart logs status clean test test-coverage help

.DEFAULT_GOAL := help

# =============================================================================
# 🚀 MAIN COMMANDS
# =============================================================================

start: _network ## Start gateway-service container
	@echo "🌐 Starting Gateway Service..."
	@cd docker && docker-compose up -d --build
	@echo "✅ Started on port 8082"

stop: ## Stop gateway-service container
	@echo "🛑 Stopping Gateway Service..."
	@cd docker && docker-compose down
	@echo "✅ Stopped"

restart: stop start ## Restart container

logs: ## View logs
	@cd docker && docker-compose logs -f

status: ## Show container status
	@cd docker && docker-compose ps

clean: ## Remove container and images
	@cd docker && docker-compose down -v --rmi local 2>/dev/null || true
	@echo "✅ Gateway Service cleaned!"

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	@go test -cover ./...

help: ## Show this help
	@echo "🌐 Gateway Service"
	@echo ""
	@echo "Usage: make [command]"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-12s %s\n", $$1, $$2}'
	@echo ""
	@echo "Quick Start: make start"

# =============================================================================
# 🔧 INTERNAL
# =============================================================================

_network:
	@docker network create docker_barrest_network 2>/dev/null || true
