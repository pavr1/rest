# Bar-Restaurant Root Makefile
# Orchestrates all services

.PHONY: test test-data test-session test-gateway test-menu test-coverage start stop restart status logs clean fresh help

.DEFAULT_GOAL := help

# =============================================================================
# 🧪 TESTING
# =============================================================================

test: ## Run all tests across all services
	@echo "🧪 Running all tests..."
	@echo ""
	@echo "📦 Data Service Tests"
	@echo "─────────────────────"
	@cd data-service && go test -v ./...
	@echo ""
	@echo "🔐 Session Service Tests"
	@echo "────────────────────────"
	@cd session-service && go test -v ./...
	@echo ""
	@echo "🌐 Gateway Service Tests"
	@echo "────────────────────────"
	@cd gateway-service && go test -v ./...
	@echo ""
	@echo "🍽️ Menu Service Tests"
	@echo "─────────────────────"
	@cd menu-service && go test -v ./...
	@echo ""
	@echo "✅ All tests complete!"

test-data: ## Run data-service tests only
	@echo "📦 Running Data Service tests..."
	@cd data-service && go test -v ./...

test-session: ## Run session-service tests only
	@echo "🔐 Running Session Service tests..."
	@cd session-service && go test -v ./...

test-gateway: ## Run gateway-service tests only
	@echo "🌐 Running Gateway Service tests..."
	@cd gateway-service && go test -v ./...

test-menu: ## Run menu-service tests only
	@echo "🍽️ Running Menu Service tests..."
	@cd menu-service && go test -v ./...

test-coverage: ## Run all tests with coverage
	@echo "🧪 Running tests with coverage..."
	@cd data-service && go test -cover ./...
	@cd session-service && go test -cover ./...
	@cd gateway-service && go test -cover ./...
	@cd menu-service && go test -cover ./...

# =============================================================================
# 🚀 SERVICE MANAGEMENT
# =============================================================================

start: ## Start all services
	@echo "🍺 Starting all Bar-Restaurant services..."
	@cd data-service && make start
	@echo "⏳ Waiting for database..."
	@sleep 3
	@cd session-service && make start
	@sleep 2
	@cd menu-service && make start
	@sleep 2
	@cd gateway-service && make start
	@sleep 2
	@cd ui-service && make start
	@echo "✅ All services started!"
	@echo ""
	@echo "┌─────────────────────────────────────────────────────────────┐"
	@echo "│                   🍺 Bar-Restaurant                        │"
	@echo "├─────────────────────────────────────────────────────────────┤"
	@echo "│  🌐 UI:        http://localhost:3000                       │"
	@echo "│  🐘 pgAdmin:   http://localhost:8080                       │"
	@echo "│                (admin@barrest.com / admin123)              │"
	@echo "│  🐳 Portainer: http://localhost:9000                       │"
	@echo "└─────────────────────────────────────────────────────────────┘"

stop: ## Stop all services
	@echo "🛑 Stopping all services..."
	@cd ui-service && make stop
	@cd gateway-service && make stop
	@cd menu-service && make stop
	@cd session-service && make stop
	@cd data-service && make stop
	@echo "✅ All services stopped!"

restart: stop start ## Restart all services

logs: ## View logs for a service (usage: make logs s=gateway)
	@if [ -z "$(s)" ]; then \
		echo "Usage: make logs s=<service>"; \
		echo "  Services: data, session, gateway, menu, ui"; \
	else \
		cd $(s)-service && make logs; \
	fi

status: ## Show status of all services
	@echo "📊 Service Status"
	@echo ""
	@echo "📦 Data Service:"
	@cd data-service && make status
	@echo ""
	@echo "🔐 Session Service:"
	@cd session-service && make status
	@echo ""
	@echo "🍽️ Menu Service:"
	@cd menu-service && make status
	@echo ""
	@echo "🌐 Gateway Service:"
	@cd gateway-service && make status
	@echo ""
	@echo "🎨 UI Service:"
	@cd ui-service && make status

clean: ## Clean all services
	@echo "🧹 Cleaning all services..."
	@cd ui-service && make clean
	@cd gateway-service && make clean
	@cd menu-service && make clean
	@cd session-service && make clean
	@cd data-service && make clean
	@echo "✅ All cleaned!"

fresh: clean ## Fresh install of all services
	@echo "🍺 Fresh install of all services..."
	@cd data-service && make fresh
	@sleep 2
	@cd session-service && make start
	@sleep 2
	@cd menu-service && make start
	@sleep 2
	@cd gateway-service && make start
	@sleep 2
	@cd ui-service && make start
	@echo "✅ Fresh install complete!"
	@echo ""
	@echo "┌─────────────────────────────────────────────────────────────┐"
	@echo "│                   🍺 Bar-Restaurant                        │"
	@echo "├─────────────────────────────────────────────────────────────┤"
	@echo "│  🌐 UI:        http://localhost:3000                       │"
	@echo "│  🐘 pgAdmin:   http://localhost:8080                       │"
	@echo "│                (admin@barrest.com / admin123)              │"
	@echo "│  🐳 Portainer: http://localhost:9000                       │"
	@echo "└─────────────────────────────────────────────────────────────┘"

# =============================================================================
# 📋 HELP
# =============================================================================

help: ## Show this help
	@echo "🍺 Bar-Restaurant Application"
	@echo ""
	@echo "Usage: make [command]"
	@echo ""
	@echo "Testing:"
	@grep -E '^test[a-zA-Z_-]*:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-18s %s\n", $$1, $$2}'
	@echo ""
	@echo "Services:"
	@grep -E '^(start|stop|restart|status|logs|clean|fresh):.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-18s %s\n", $$1, $$2}'
	@echo ""
	@echo "Quick Start:"
	@echo "  make fresh           # Clean install everything"
	@echo "  make test            # Run all tests"
	@echo "  make logs s=gateway  # View gateway service logs"
	@echo "  make logs s=menu     # View menu service logs"
