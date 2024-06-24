.PHONY: help
help: ## Show help messages.
	@grep -E '^[0-9a-zA-Z\/_-]+:(.*?## .*)?$$' $(MAKEFILE_LIST) | sed 's/^[^:]*Makefile://' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: lint
lint: ## Lint it
	golangci-lint run --verbose ./...
