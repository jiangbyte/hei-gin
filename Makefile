.PHONY: dev build run test lint clean scaffold list-imports

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
	swag init --parseDependency --parseInternal --parseDepth 3 && \
	cp docs/swagger.json sdk/app/swagger.json

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
test:         ## Run tests
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
	rm -f sdk/app/swagger.json

help:         ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
