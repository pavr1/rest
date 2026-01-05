# Bar Restaurant - UI Service Makefile

.PHONY: start stop restart logs status clean test help

.DEFAULT_GOAL := help

# =============================================================================
# 🚀 DOCKER COMMANDS
# =============================================================================

start: ## Start UI service container
	@echo "🎨 Starting UI Service..."
	@cd docker && docker-compose up -d --build
	@echo "✅ UI Service started on http://localhost:3000"

stop: ## Stop UI service container
	@echo "🛑 Stopping UI Service..."
	@cd docker && docker-compose down
	@echo "✅ UI Service stopped!"

restart: stop start ## Restart UI service

logs: ## View UI service logs
	@cd docker && docker-compose logs -f

status: ## Show UI service status
	@echo "📊 UI Service Status:"
	@cd docker && docker-compose ps

clean: ## Remove container and images
	@cd docker && docker-compose down -v --rmi local 2>/dev/null || true
	@echo "✅ UI Service cleaned!"

# =============================================================================
# 🧪 TESTING & VALIDATION
# =============================================================================

test: ## Run UI health check
	@echo "🧪 Testing UI Service..."
	@curl -sf http://localhost:3000/health > /dev/null && echo "✅ UI Service is healthy!" || echo "❌ UI Service is not responding"

# =============================================================================
# 💻 LOCAL DEVELOPMENT (optional, requires npm)
# =============================================================================

install: ## Install npm dependencies (for local dev)
	@echo "📦 Installing dependencies..."
	@npm install
	@echo "✅ Dependencies installed!"

dev: ## Start local development server
	@echo "🚀 Starting local development server..."
	@npx http-server . -p 3000 -c-1

# =============================================================================
# 📋 HELP
# =============================================================================

help: ## Show this help
	@echo "🎨 Bar Restaurant UI Service"
	@echo ""
	@echo "Usage: make [command]"
	@echo ""
	@echo "Docker Commands:"
	@grep -E '^(start|stop|restart|logs|status|clean):.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-12s %s\n", $$1, $$2}'
	@echo ""
	@echo "Testing:"
	@grep -E '^test:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-12s %s\n", $$1, $$2}'
	@echo ""
	@echo "Local Development:"
	@grep -E '^(install|dev):.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-12s %s\n", $$1, $$2}'

