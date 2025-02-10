# kubectl-vault-login
Ever needed short-lived and fine-grained `kubectl` access to `Kubernetes`-Cluster for some operations during CI/CD?

Well, `kubectl-vault-login` allows you to do exactly this!

By leveraging [HashiCorp Vaults Kubernetes Secrets Engine](https://developer.hashicorp.com/vault/docs/secrets/kubernetes) you can create `Service Accounts` and `Service Account tokens` with a **tight RBAC** and a **low TTL** - making it incredibly useful for CI/CD Operations, such as `kubectl apply` commands.

## How does it work
The [HashiCorp Vaults Kubernetes Secrets Engine](https://developer.hashicorp.com/vault/docs/secrets/kubernetes) can operate in 3 modes:

1. [Creating a `Service Account Token` for an already existing `Service Account` with a pre-existing `Role` & `Rolebinding`](./mode-01.md)
2. [Creating a `Service Account` and a `Service Account Token` with a pre-existing `Role` & `Rolebinding`](./mode-02.md)
3. [Creating a `Service Account`, a `Service Account Token` and the `Rolebinding` for an pre-existing `Role`](./mode-03.md)

Every resource created by `Vault` will automatically revoked after the lease is expired.

## Getting started
For every mode, the steps are the same:

1. Install the plugin
2. Configure Kubernetes to allow `Vault` creating the resources
3. Configure [HashiCorp Vaults Kubernetes Secrets Engine](https://developer.hashicorp.com/vault/docs/secrets/kubernetes)
4. Patch your `kubeconfig` to use `kubectl-vault-login`:

```yaml
# $KUBECONFIG
[...]
users:
- name: kind-kind
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      command: kubectl-vault-login
      args:
        - --role=kind
```


5. Run any `kubectl` plugin that is allowed in your RBAC-setup

**Checkout the [Guides](./mode-01.md)**
