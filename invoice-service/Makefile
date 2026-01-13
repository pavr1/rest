.PHONY: start stop restart logs status clean test test-coverage

# Invoice Service
# =============================================================================

start: ## Start invoice-service container
	@echo "📄 Starting Invoice Service..."
	@docker-compose -f docker/docker-compose.yml up -d
	@echo "✅ Started on port 8092"

stop: ## Stop invoice-service container
	@echo "🛑 Stopping Invoice Service..."
	@docker-compose -f docker/docker-compose.yml down
	@echo "✅ Stopped"

restart: ## Restart invoice-service container
	@echo "🔄 Restarting Invoice Service..."
	@docker-compose -f docker/docker-compose.yml down
	@docker-compose -f docker/docker-compose.yml up -d
	@echo "✅ Restarted"

logs: ## View invoice-service logs
	@docker logs -f barrest_invoice_service

status: ## Show invoice-service container status
	@docker ps --filter "name=barrest_invoice_service" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

clean: ## Remove invoice-service container and images
	@docker-compose -f docker/docker-compose.yml down -v --remove-orphans
	@docker rmi barrest_invoice_service 2>/dev/null || true
	@echo "✅ Invoice Service cleaned!"

test: ## Run invoice-service tests
	@echo "🧪 Running Invoice Service tests..."
	@go test -v ./...

test-coverage: ## Run invoice-service tests with coverage
	@echo "🧪 Running Invoice Service tests with coverage..."
	@go test -v -cover ./...

help: ## Show this help
	@echo "📄 Invoice Service"
	@echo ""
	@echo "Usage: make [command]"
	@echo ""
	@echo "Commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
