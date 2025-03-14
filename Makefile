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
kind: ## start local kind cluster with vault SA configured
	kind create cluster --config=scripts/kind-config.yaml

.PHONY: vault
vault:	## start local vault server with k8s secret engine configured
	vault server \
		-dev \
		-dev-listen-address=0.0.0.0:8200 \
		-dev-root-token-id=root &

.PHONY: setup-mode-01
setup-mode-01: ## setup vault in mode 01
	kubectl apply -f ./scripts/mode-01
	./scripts/mode-01/setup-vault.sh

.PHONY: setup-mode-02
setup-mode-02: ## setup vault in mode 02
	kubectl apply -f ./scripts/mode-02
	./scripts/mode-02/setup-vault.sh

.PHONY: setup-mode-03
setup-mode-03: ## setup vault in mode 03
	kubectl apply -f ./scripts/mode-03
	./scripts/mode-03/setup-vault.sh

.PHONY: teardown
teardown:
	kind delete cluster
	kill -9 $(shell pgrep -x vault)
