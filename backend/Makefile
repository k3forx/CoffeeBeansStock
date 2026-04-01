DB_URL ?= postgres://coffee_user:coffee_pass@localhost:5432/coffee_stock?sslmode=disable
MIGRATIONS_DIR = migrations

.PHONY: help run dev sqlc-generate migrate-up migrate-down test build clean lint generate-api

help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: ## Run the server
	go run cmd/server/main.go

dev: ## Run the server with hot reload (air)
	$(GOBIN)/air

sqlc-generate: ## Generate sqlc code
	$(GOBIN)/sqlc generate

migrate-up: ## Run migrations up
	$(GOBIN)/migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down: ## Run migrations down (1 step)
	$(GOBIN)/migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

test: ## Run tests
	go test -v -race ./...

build: ## Build the binary
	go build -o bin/server cmd/server/main.go

clean: ## Clean build artifacts
	rm -rf bin/

GOLANGCI_LINT_VERSION ?= v2.1.6
TARGETS ?= ./...

generate-api: ## Generate Go types from OpenAPI spec
	$(GOBIN)/oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml

lint: ## Run golangci-lint (use TARGETS="./path/..." to specify targets)
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION); \
	fi
	golangci-lint run $(TARGETS)
