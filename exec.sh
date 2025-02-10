#!/usr/bin/env bash
cat <<EOF
{
  "apiVersion": "client.authentication.k8s.io/v1beta1",
  "kind": "ExecCredential",
  "status": {
    "token": "$(vault write -format=json -field=service_account_token kubernetes/creds/"${VAULT_ROLE}" ttl="${TTL}" | jq -r .)"
  }
}
EOF
