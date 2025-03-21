# kubectl-vault-login
Ever needed short-lived and fine-grained `kubectl` access to `Kubernetes`-Cluster during CI/CD?

Well, `kubectl-vault-login` allows you to do exactly this!

By leveraging [HashiCorp Vaults Kubernetes Secrets Engine](https://developer.hashicorp.com/vault/docs/secrets/kubernetes) you can create `Service Accounts` and `Service Account tokens` with a **tight RBAC** and a **low TTL** - making it powerful for CI/CD Operations, such as `kubectl apply` commands.

## How does it work
The [HashiCorp Vaults Kubernetes Secrets Engine](https://developer.hashicorp.com/vault/docs/secrets/kubernetes) can operate in 3 modes:

1. [Create a ServiceAccount Token for a ServiceAccount with Role & RoleBinding](https://falcosuessgott.github.io/kubectl-vault-login/mode-01/)
2. [Create a ServiceAccount, Token and RoleBinding for a (Cluster)-Role (e.g `cluster-admin`)](https://falcosuessgott.github.io/kubectl-vault-login/mode-02/)
3. [Create a ServiceAccount, a Token, Role and RoleBinding](https://falcosuessgott.github.io/kubectl-vault-login/mode-03/)

Every resource created by `Vault` will automatically revoked once the lease is expired (minimum `600s`).

> [!TIP]
> `kubectl-vault-login` will cache the token to `~/.kube/cache/vault-login/token` (change with `$KUBECACHEDIR`) and re-use the token until expiry

![img](./docs/assets/workflow.svg)

## Getting started
For every mode, the steps are the same:

1. Install the plugin
2. Configure a Kubernetes ServiceAccount that is being used by Vault to create RBAC resources
3. Configure [HashiCorp Vaults Kubernetes Secrets Engine](https://developer.hashicorp.com/vault/docs/secrets/kubernetes)
4. Create the necessary (Cluster)-Roles and (Cluster)-RoleBindings for which the ServiceAccounts are going to be created
5. Patch your `$KUBECONFIG` to use `kubectl-vault-login` as an [`ExecCredential`](https://kubernetes.io/docs/reference/config-api/client-authentication.v1beta1/):

```bash
> kubectl config set-credentials vault \
  --exec-interactive-mode=Never \
  --exec-api-version=client.authentication.k8s.io/v1 \
  --exec-command=kubectl \
  --exec-arg=vault-login \
  --exec-arg=--role=kind # change to your role
```

```yaml
# $KUBECONFIG
[...]
users:
- name: kind-kind
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      command: kubectl
      args:
        - vault-login
        - --role=kind
```

6. Run any `kubectl` plugin that is allowed in your RBAC-setup

**Checkout the [Guides](https://falcosuessgott.github.io/kubectl-vault-login/mode-01/)**
