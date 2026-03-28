.PHONY: help run sqlc-generate migrate-up migrate-down test build clean

help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: ## Run the server
	go run cmd/server/main.go

sqlc-generate: ## Generate sqlc code
	sqlc generate

migrate-up: ## Run migrations up
	# TODO: マイグレーションツール導入後に実装

migrate-down: ## Run migrations down
	# TODO: マイグレーションツール導入後に実装

test: ## Run tests
	go test -v ./...

build: ## Build the binary
	go build -o bin/server cmd/server/main.go

clean: ## Clean build artifacts
	rm -rf bin/
