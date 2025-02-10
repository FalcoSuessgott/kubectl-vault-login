# Creating a `Service Account Token` for an already existing `Service Account` with a pre-existing `Role` & `Rolebinding`
This guide will walk you through setting up `Kubernetes` and `Vault` and its `Kubernetes Secret Engine` to create a `Service Account Token` for an already existing `Service Account` with a pre-existing `Role` & `Rolebinding`

## Prerequisites
You will need the following tools to be installed:

- [`kind`](https://kind.sigs.k8s.io)
- [`Vault`](https://developer.hashicorp.com/vault/docs/install)
- [`kubectl-vault-login`](https://falcosuessgott.github.io/kubectl-vault-login/)

## Setup `kind`
```bash
kind create cluster --config=kind-config.yml
```

you should now be able to run `kubectl commands:

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
The following manifest, creates a `Service Account` `vault-auth` and binds the role `service-account-token-creator` to it, which allows to create `Service Account Tokens`. This `Service Account` is being used by `Vault` to create `Service Account Tokens` for the `tmp-sa` `Service Account` that we will create in the next section:

```yaml
cat <<EOF | kubectl create -f -
{!../scripts/mode-01/vault-auth.yml!}
EOF
```

## Create a `Service Account` for which `Vault` creates the `Service Account Token`
This manifest creates a `Service Account` `tmp-sa` that is bound to the `role-list-pods` role that **only** allows to **list pods in the `default` namespace**:

```yaml
cat <<EOF | kubectl create -f -
{!../scripts/mode-01/tmp-sa.yml!}
EOF
```

## Configure `Vault`
Lastly, we will need to start and configure a local `Vault Server`:

```bash
vault server \
	-dev \
	dev-listen-address=0.0.0.0:8200 \
	dev-root-token-id=root
```

Now, we will configure the `Kubernetes Secrets Engine` to connect to the local `kind` Cluster with the `vault-auth` `Service Account` and creating a role `kind` that will create a `Service Account Token` for the `tmp-sa` Service Account:

```bash
{!../.envrc!}
```

```bash
{!../scripts/mode-01/setup-vault.sh!}
```

## Putting it together
Write `kind`s `kubeconfig` to a file:

```bash
kind get kubeconfig > kind-kubeconfig.yml
```

and update it, to use `kubectl-vault-login` for obtaining access:

```yaml
# kind-kubeconfig.yml
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

You will still need to be authenticated to `Vault`:

```bash
{!../.envrc!}
```

```bash
# create a pod to see some results
kubectl run nginx --image=nginx
# use the updated kubeconfig to list pods in the default namespace
KUBECONFIG=kubeconfig.yml kubectl get pods
```

And you should see:
```bash
NAME    READY   STATUS    RESTARTS   AGE
nginx   1/1     Running   0          73s
```

The role `role-list-pods` allows listing pods for the `default` namespace, but not for `kube-system`:

```bash
KUBECONFIG=vault-kubeconfig.yml k get pod -n kube-config
Error from server (Forbidden): pods is forbidden: User "system:serviceaccount:default:v-token-kind-1739680669-u5x0uqreffqt8hf2qdydpksf" cannot list resource "pods" in API group "" in the namespace "kube-system"
```

## Teardown
Tear everything down by running:
```bash
kind delete cluster
kill -9 $(pgrep -x vault)
```
