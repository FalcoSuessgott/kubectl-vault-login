# Usage

!!! tip
    All of Vaults [Environment Variables](https://developer.hashicorp.com/vault/docs/commands) are supported.

```bash
> kubectl-vault-login -h
A kubectl plugin to to obtain access to a kubernetes cluster via HashiCorp Vaults Kubernetes secrets engine

Usage:
  kubectl-vault-login [flags]

Flags:
  -a, --audiences string   A comma separated string containing the intended audiences of the generated Kubernetes service account (VAULT_K8S_LOGIN_AUDIENCES)
  -c, --crb                If true, generate a ClusterRoleBinding to grant permissions across the whole cluster instead of within a namespace (VAULT_K8S_LOGIN_CRB)
  -h, --help               help for kubectl-vault-login
  -m, --mount string       The Kubernetes secrets mount path (VAULT_K8S_LOGIN_MOUNT) (default "kubernetes")
  -n, --ns string          The name of the Kubernetes namespace in which to generate the credentials (VAULT_K8S_LOGIN_NAMESPACE)
  -r, --role string        The name of the role to generate credentials for (VAULT_K8S_LOGIN_ROLE)
  -t, --ttl string         The ttl of the generated Kubernetes service account (VAULT_K8S_LOGIN_TTL) (default "1h")
```
