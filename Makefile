default: help

.PHONY: help
help: ## list makefile targets
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

PHONY: fmt
fmt: ## format go files
	gofumpt -w .
	gci write .
	pre-commit run -a

.PHONY: docs
docs: ## render docs locally
	mkdocs serve

PHONY: test
test: ## test
	gotestsum -- -v --shuffle=on -race -coverprofile="coverage.out" -covermode=atomic ./...

PHONY: lint
lint: ## lint go files
	golangci-lint run -c .golang-ci.yml

.PHONY: kind
kind:
	kind create cluster --config=kind-config.yml

.PHONY: vault
vault:
	vault server \
    	-dev \
    	-dev-listen-address=0.0.0.0:8200 \
    	-dev-root-token-id=root &

.PHONY: teardown
teardown:
	kind delete cluster
	kill -9 $(pgrep -x vault)
