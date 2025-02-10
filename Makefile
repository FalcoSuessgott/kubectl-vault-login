default: help

.PHONY: help
help: ## list makefile targets
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

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
