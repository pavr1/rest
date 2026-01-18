# Bar-Restaurant Root Makefile
# Orchestrates all services

.PHONY: test test-data test-session test-gateway test-menu test-inventory test-invoice test-orders test-coverage start stop restart status logs clean fresh help

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
	@echo "🧾 Invoice Service Tests"
	@echo "────────────────────────"
	@cd invoice-service && go test -v ./...
	@echo ""
	@# Future services (uncomment when implemented):
	@# @echo "📦 Inventory Service Tests"
	@# @cd inventory-service && go test -v ./...
	@# @echo "📋 Orders Service Tests"
	@# @cd orders-service && go test -v ./...
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

test-inventory: ## Run inventory-service tests only (future)
	@echo "📦 Running Inventory Service tests..."
	@if [ -d "inventory-service" ]; then cd inventory-service && go test -v ./...; else echo "⚠️ inventory-service not yet implemented"; fi

test-invoice: ## Run invoice-service tests only
	@echo "🧾 Running Invoice Service tests..."
	@cd invoice-service && go test -v ./...

test-orders: ## Run orders-service tests only (future)
	@echo "📋 Running Orders Service tests..."
	@if [ -d "orders-service" ]; then cd orders-service && go test -v ./...; else echo "⚠️ orders-service not yet implemented"; fi

test-coverage: ## Run all tests with coverage
	@echo "🧪 Running tests with coverage..."
	@cd data-service && go test -cover ./...
	@cd session-service && go test -cover ./...
	@cd gateway-service && go test -cover ./...
	@cd menu-service && go test -cover ./...
	@cd invoice-service && go test -cover ./...
	@# Future services (uncomment when implemented):
	@# @cd inventory-service && go test -cover ./...
	@# @cd orders-service && go test -cover ./...

# =============================================================================
# 🚀 SERVICE MANAGEMENT
# =============================================================================

start: ## Start all services
	@echo "🍺 Starting all Bar-Restaurant services..."
	@echo ""
	@echo "Level 0-1: Database + Data Service"
	@cd data-service && make start
	@echo "⏳ Waiting for database..."
	@sleep 3
	@echo ""
	@echo "Level 2: Auth + Inventory"
	@cd session-service && make start
	@# @cd inventory-service && make start  # Future
	@sleep 2
	@echo ""
	@echo "Level 3: Business Services"
	@cd invoice-service && make start
	@cd menu-service && make start
	@# @cd customer-service && make start  # Future
	@# @cd karaoke-service && make start  # Future
	@# @cd promotion-service && make start  # Future
	@sleep 2
	@echo ""
	@echo "Level 4: Orders"
	@# @cd orders-service && make start  # Future
	@sleep 1
	@echo ""
	@echo "Level 5: Gateway"
	@cd gateway-service && make start
	@sleep 2
	@echo ""
	@echo "Level 6: UI"
	@cd ui-service && make start
	@echo "✅ All services started!"
	@echo ""
	@echo "┌─────────────────────────────────────────────────────────────┐"
	@echo "│                   🍺 Bar-Restaurant                         │"
	@echo "├─────────────────────────────────────────────────────────────┤"
	@echo "│  🌐 UI:        http://localhost:3000                        │"
	@echo "│  🐘 pgAdmin:   http://localhost:8080                        │"
	@echo "│                (admin@barrest.com / admin123)               │"
	@echo "│  🐳 Portainer: http://localhost:9000                        │"
	@echo "└─────────────────────────────────────────────────────────────┘"

stop: ## Stop all services (reverse order)
	@echo "🛑 Stopping all services..."
	@echo "Level 6: UI"
	@cd ui-service && make stop
	@echo "Level 5: Gateway"
	@cd gateway-service && make stop
	@echo "Level 4: Orders"
	@# @cd orders-service && make stop  # Future
	@echo "Level 3: Business Services"
	@# @cd promotion-service && make stop  # Future
	@# @cd karaoke-service && make stop  # Future
	@# @cd customer-service && make stop  # Future
	@cd menu-service && make stop
	@cd invoice-service && make stop
	@echo "Level 2: Auth + Inventory"
	@# @cd inventory-service && make stop  # Future
	@cd session-service && make stop
	@echo "Level 0-1: Data + Database"
	@cd data-service && make stop
	@echo "✅ All services stopped!"

restart: stop start ## Restart all services

logs: ## View logs for a service (usage: make logs s=gateway)
	@if [ -z "$(s)" ]; then \
		echo "Usage: make logs s=<service>"; \
		echo "  Services: data, session, gateway, menu, invoice, ui"; \
		echo "  Future: inventory, orders, customer, karaoke, promotion"; \
	else \
		cd $(s)-service && make logs; \
	fi

status: ## Show status of all services
	@echo "📊 Service Status"
	@echo ""
	@echo "═══════════════════════════════════════════════════════════════"
	@echo "Level 0-1: Infrastructure"
	@echo "═══════════════════════════════════════════════════════════════"
	@echo "📦 Data Service (8086):"
	@cd data-service && make status
	@echo ""
	@echo "═══════════════════════════════════════════════════════════════"
	@echo "Level 2: Auth + Inventory"
	@echo "═══════════════════════════════════════════════════════════════"
	@echo "🔐 Session Service (8087):"
	@cd session-service && make status
	@echo ""
	@echo "📦 Inventory Service (8090): [NOT IMPLEMENTED]"
	@# @cd inventory-service && make status  # Future
	@echo ""
	@echo "═══════════════════════════════════════════════════════════════"
	@echo "Level 3: Business Services"
	@echo "═══════════════════════════════════════════════════════════════"
	@echo "🧾 Invoice Service (8092):"
	@cd invoice-service && make status
	@echo ""
	@echo "🍽️ Menu Service (8088):"
	@cd menu-service && make status
	@echo ""
	@echo "👥 Customer Service (8095): [NOT IMPLEMENTED]"
	@# @cd customer-service && make status  # Future
	@echo ""
	@echo "🎤 Karaoke Service (8093): [NOT IMPLEMENTED]"
	@# @cd karaoke-service && make status  # Future
	@echo ""
	@echo "🎁 Promotion Service (8094): [NOT IMPLEMENTED]"
	@# @cd promotion-service && make status  # Future
	@echo ""
	@echo "═══════════════════════════════════════════════════════════════"
	@echo "Level 4: Orders"
	@echo "═══════════════════════════════════════════════════════════════"
	@echo "📋 Orders Service (8089): [NOT IMPLEMENTED]"
	@# @cd orders-service && make status  # Future
	@echo ""
	@echo "═══════════════════════════════════════════════════════════════"
	@echo "Level 5-6: Gateway + UI"
	@echo "═══════════════════════════════════════════════════════════════"
	@echo "🌐 Gateway Service (8082):"
	@cd gateway-service && make status
	@echo ""
	@echo "🎨 UI Service (3000):"
	@cd ui-service && make status

clean: ## Clean all services (reverse order)
	@echo "🧹 Cleaning all services..."
	@cd ui-service && make clean
	@cd gateway-service && make clean
	@# @cd orders-service && make clean  # Future
	@# @cd promotion-service && make clean  # Future
	@# @cd karaoke-service && make clean  # Future
	@# @cd customer-service && make clean  # Future
	@cd menu-service && make clean
	@cd invoice-service && make clean
	@cd inventory-service && make clean 
	@cd session-service && make clean
	@cd data-service && make clean
	@echo "✅ All cleaned!"

fresh: clean ## Fresh install of all services
	@echo "🍺 Fresh install of all services..."
	@echo ""
	@echo "Level 0-1: Database + Data Service"
	@cd data-service && make fresh
	@sleep 2
	@echo ""
	@echo "Level 2: Auth + Inventory"
	@cd session-service && make start
	@cd inventory-service && make start
	@sleep 2
	@echo ""
	@echo "Level 3: Business Services"
	@cd invoice-service && make start
	@cd menu-service && make start
	@# @cd customer-service && make start  # Future
	@# @cd karaoke-service && make start  # Future
	@# @cd promotion-service && make start  # Future
	@sleep 2
	@echo ""
	@echo "Level 4: Orders"
	@# @cd orders-service && make start  # Future
	@sleep 1
	@echo ""
	@echo "Level 5: Gateway"
	@cd gateway-service && make start
	@sleep 2
	@echo ""
	@echo "Level 6: UI"
	@cd ui-service && make start
	@echo "✅ Fresh install complete!"
	@echo ""
	@echo "┌─────────────────────────────────────────────────────────────┐"
	@echo "│                   🍺 Bar-Restaurant                         │"
	@echo "├─────────────────────────────────────────────────────────────┤"
	@echo "│  🌐 UI:        http://localhost:3000                        │"
	@echo "│  🐘 pgAdmin:   http://localhost:8080                        │"
	@echo "│                (admin@barrest.com / admin123)               │"
	@echo "│  🐳 Portainer: http://localhost:9000                        │"
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
	@echo "Service Startup Order (by level):"
	@echo "  Level 0-1: data-service (8086)"
	@echo "  Level 2:   session-service (8087), inventory-service (8090)*"
	@echo "  Level 3:   invoice-service (8092), menu-service (8088),"
	@echo "             customer-service (8095)*, karaoke-service (8093)*,"
	@echo "             promotion-service (8094)*"
	@echo "  Level 4:   orders-service (8089)*"
	@echo "  Level 5:   gateway-service (8082)"
	@echo "  Level 6:   ui-service (3000)"
	@echo "  * = Not yet implemented"
	@echo ""
	@echo "Quick Start:"
	@echo "  make fresh           # Clean install everything"
	@echo "  make test            # Run all tests"
	@echo "  make logs s=gateway  # View gateway service logs"
	@echo "  make logs s=menu     # View menu service logs"
