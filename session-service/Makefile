.PHONY: build run test clean deps help

.DEFAULT_GOAL := help

build: deps ## Build the service
	@echo "🔨 Building session-service..."
	@go build -o main .
	@echo "✅ Build complete!"

run: build ## Run locally
	@echo "🔐 Starting session-service..."
	@./main

deps: ## Download dependencies
	@go mod tidy
	@go mod download

test: ## Run tests
	@go test -v ./...

clean: ## Clean build artifacts
	@rm -f main
	@echo "✅ Cleaned!"

help: ## Show help
	@echo "🔐 Session Service"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-10s %s\n", $$1, $$2}'
