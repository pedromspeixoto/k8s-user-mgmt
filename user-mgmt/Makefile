# General purpose targets

.PHONY: help
help: ## Display available commands.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Docker compose targets.

.PHONY: build
build: ## Build the application
	./scripts/build.sh

.PHONY: run
run: ## Run the application
	go run cmd/main.go