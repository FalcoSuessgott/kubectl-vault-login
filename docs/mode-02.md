# Create a ServiceAccount, Token and RoleBinding for a (Cluster)-Role (e.g `cluster-admin`)
!!! tip
    This guide will walk you through setting up `kind` and `Vault` and its Kubernetes Secret Engine to create a Service Account, Token and RoleBinding for the predefined `cluster-admin` ClusterRole

!!! warning
    The `cluster-admin` role can do anything in every namespace.

    **Use with caution**

## Prerequisites
You will need the following tools to be installed:

- [`kind`](https://kind.sigs.k8s.io)
- [`Vault`](https://developer.hashicorp.com/vault/docs/install)
- [`kubectl-vault-login`](https://falcosuessgott.github.io/kubectl-vault-login/)

## Setup `kind`
```bash
cat <<EOF >>kind-config.yaml
{!../scripts/kind-config.yaml!}
EOF
kind create cluster --config=kind-config.yaml
```

you should now be able to run `kubectl` commands:

```bash
kubectl get ns
NAME                 STATUS   AGE
default              Active   64m
kube-node-lease      Active   64m
kube-public          Active   64m
kube-system          Active   64m
local-path-storage   Active   63m
```

## Configure `Vault` access
The following manifest, creates a ServiceAccount `vault-auth` and assigns it the role `cluster-admin-creator`, which allows to create Service Account, Tokens assigning them the a (Cluster)-Role.

!!! tip
    This Service Account is going to be used by `Vault`

!!! note
    **Kubernetes prevents users (including service accounts) from granting RBAC permissions they do not already have themselves.
    Thats why we have to assign `bind` and `escalate` as `verbs` for `clusterrolebindings`.

```yaml
cat <<EOF | kubectl apply -f -
{!../scripts/mode-02/vault-auth.yml!}
EOF
```

## Configure `Vault`
Lastly, we will need to start and configure a local `Vault Server`:

```bash
vault server \
	-dev \
	-dev-listen-address=0.0.0.0:8200 \
	-dev-root-token-id=root
```

Authenticate to `Vault` and check with `vault status`:

```bash
{!../.envrc!}
> vault status
Key             Value
---             -----
Seal Type       shamir
Initialized     true
Sealed          false
Total Shares    1
Threshold       1
Version         1.18.3
Build Date      2024-12-16T14:00:53Z
Storage Type    inmem
Cluster Name    vault-cluster-4cab3957
Cluster ID      597257da-8e8d-6147-c379-e93e3a6013c7
HA Enabled      false
```

Now, we will configure the Kubernetes Secrets Engine to connect to the local `kind` Cluster with the `vault-auth` ServiceAccount and create a role `kind` that will create the ServiceAcccount, Token and RoleBinding:

!!! important
    Note the `kubernetes_role_type` and `kubernetes_role_name`

```bash
{!../scripts/mode-02/setup-vault.sh!}
```

## Putting it together
Write `kind`s `kubeconfig` to a file:

```bash
kind get kubeconfig > kubeconfig.yml
```

and update it, to use `kubectl-vault-login` for authentication:

```bash
KUBECONFIG=./kubeconfig.yml kubectl config set-credentials vault \
  --exec-interactive-mode=Never \
  --exec-api-version=client.authentication.k8s.io/v1 \
  --exec-command=kubectl \
  --exec-arg=vault \
  --exec-arg=login \
  --exec-arg=--role=kind
```

```bash
> cat kubeconfig.yml
[...]
users:
- name: vault
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1
      args:
      - vault
      - login
      - --role=kind
      command: kubectl
      env: null
      interactiveMode: Never
      provideClusterInfo: false
```

```bash
# check SA has been created
> KUBECONFIG=kubeconfig.yml kubectl get sa
NAME                                               SECRETS   AGE
v-token-kind-1739829804-zbmswmhaet1qelxccyo97uux   0         25s
# check clusterrrolebinding was created
> KUBECONFIG=kubeconfig.yml kubectl get clusterrolebindings -n default
NAME                                             ROLE                      AGE
v-token-kind-1739829804-zbmswmhaet1qelxccyo97uux ClusterRole/cluster-admin 59s
```

## Teardown
Tear everything down by running:
```bash
kind delete cluster
kill -9 $(pgrep -x vault)
```
