# Exploring the [Vault Kubernetes Secrets Engine](https://developer.hashicorp.com/vault/docs/secrets/kubernetes)

Operates in 3 modes:
1. Create a SA Token
2. Create a SA and a SA Token and bind a pre-existing role
3. Create a SA, a SA Token, a Role and a Rolebinding

Prerequisites:
- kind
- vault
- kubectl
- jq
