.PHONY: dev build run test lint clean scaffold list-imports test-root test-sdk test-plugin-sys test-plugin-client test-plugin-im test-all

.PHONY: swag swag-serve
GO_MODULES := . api sdk plugins/plugin-sys plugins/plugin-client plugins/plugin-im

# ── Development ─────────────────────────────────────────────────────
dev:          ## Start with hot-reload (requires air)
	@command -v air >/dev/null 2>&1 || { \
		echo "Installing air..."; go install github.com/air-verse/air@latest; \
	}
	air

swag:         ## Generate swagger docs
	@command -v swag >/dev/null 2>&1 || { \
		echo "Installing swag..."; go install github.com/swaggo/swag/cmd/swag@latest; \
	}
	swag init --parseDependency --parseInternal --parseDepth 3

swag-serve: swag ## Generate swagger docs & run server
	go run main.go

build:        ## Build binary
	go build -o bin/hei-gin main.go

run:          ## Build & run
	go run main.go

# ── Code Generation ─────────────────────────────────────────────────
scaffold:     ## Create new plugin: make scaffold name=plugin-xxx
	go run cmd/codegen/main.go scaffold $(name)

list-plugins: ## List all plugins
	go run cmd/codegen/main.go list

gen-imports:  ## Regenerate blank imports in main.go
	go run cmd/codegen/main.go gen-imports

# ── Database ────────────────────────────────────────────────────────
migrate:      ## Apply database migrations
	go run cmd/migrate/main.go

migrate-dry:  ## Preview database migrations (dry-run)
	go run cmd/migrate/main.go -skip-seed

# ── Quality ─────────────────────────────────────────────────────────
test: test-all ## Run tests

test-root:    ## Run tests for root module
	go test ./...

test-sdk:     ## Run tests for sdk module
	(cd sdk && go test ./...)

test-plugin-sys: ## Run tests for plugin-sys module
	(cd plugins/plugin-sys && go test ./...)

test-plugin-client: ## Run tests for plugin-client module
	(cd plugins/plugin-client && go test ./...)

test-plugin-im: ## Run tests for plugin-im module
	(cd plugins/plugin-im && go test ./...)

test-all:     ## Run tests for all workspace modules
	@for module in $(GO_MODULES); do \
		echo "==> go test $$module"; \
		(cd $$module && go test ./...) || exit 1; \
	done

lint:         ## Run linter
	@for module in $(GO_MODULES); do \
		echo "==> go vet $$module"; \
		(cd $$module && go vet ./...) || exit 1; \
	done

clean:        ## Clean build artifacts
	rm -rf .air_tmp bin/
	rm -f docs/swagger.json docs/swagger.yaml docs/docs.go

help:         ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
