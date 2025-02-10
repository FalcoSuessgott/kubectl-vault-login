#!/usr/bin/env bash
set -ex

K8S_JWT_TOKEN=$(kubectl get secret vault-auth-token -o jsonpath="{.data.token}" | base64 -d)
K8S_CA_CERT=$(kubectl get secret vault-auth-token -o jsonpath="{['data']['ca\.crt']}" | base64 -d)

vault secrets enable kubernetes
vault write -f kubernetes/config \
    kubernetes_host="https://127.0.0.1:6443" \
    kubernetes_ca_cert="$K8S_CA_CERT" \
    service_account_jwt="$K8S_JWT_TOKEN"

vault write kubernetes/roles/kind \
    allowed_kubernetes_namespaces="default" \
    kubernetes_role_name="role-list-pods" \
    token_default_ttl="10m"
