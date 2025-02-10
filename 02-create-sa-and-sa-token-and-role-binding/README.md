# 02: Generate a Kubernetes role, role binding, service account and token

## kind
```bash
kind create cluster --config=kind-config.yml
```

# SA token (vault-auth) + token ued by vault that is allowed to create service account tokens
```bash
kubectl apply -f vault-auth.yml
```

## Vault
```bash
vault server \
    -dev \
    -dev-listen-address=0.0.0.0:8200 \
    -dev-root-token-id=root
```

**Kubernetes prevents users (including service accounts) from granting RBAC permissions they do not already have themselves**

## Configure vault
```bash
K8S_JWT_TOKEN=$(kubectl get secret vault-auth-token -o jsonpath="{.data.token}" | base64 -d)
K8S_CA_CERT=$(kubectl get secret vault-auth-token -o jsonpath="{['data']['ca\.crt']}" | base64 -d)
vault secrets enable kubernetes
vault write -f kubernetes/config \
    kubernetes_host="https://127.0.0.1:6443" \
    kubernetes_ca_cert="$K8S_CA_CERT" \
    service_account_jwt="$K8S_JWT_TOKEN"
vault write kubernetes/roles/test \
    allowed_kubernetes_namespaces="default" \
    kubernetes_role_name="role-list-pods" \
    token_default_ttl="10m"

vault write \
    kubernetes/creds/test \
    ttl=20m
```

## create a pod to show we can list pods
```bash
kubectl run nginx --image=nginx
```

## fetch credentials
```bash
kubectl apply -f tmp-sa.yml
token=$(vault write -field=service_account_token \
    kubernetes/creds/test \
    kubernetes_namespace=default \
    ttl=20m)
curl -H "Authorization: Bearer $token" -sk $(kubectl config view --minify -o 'jsonpath={.clusters[].cluster.server}')/api/v1/namespaces/default/pods | jq .
```

## example kubeconfig
```
kind get kubeconfig > vault-kubeconfig.yml
```

```yaml
# vault-kubeconfig.yml
[...]
users:
- name: kind-kind
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      command: ./exec.sh
      env:
        - name: VAULT_ROLE
          value: test
        - name: K8S_NAMESPACE
          value: default
        - name: TTL
          value: 20m
```

## Example exec script:
```bash
#!/usr/bin/env bash
cat <<EOF
{
  "apiVersion": "client.authentication.k8s.io/v1beta1",
  "kind": "ExecCredential",
  "status": {
    "token": "$(vault write -format=json -field=service_account_token kubernetes/creds/"${VAULT_ROLE}" kubernetes_namespace="${K8S_NAMESPACE}" ttl="${TTL}" | jq -r .)"
  }
}
EOF
```

SA works for namespace default:
```bash
KUBECONFIG=vault-kubeconfig.yml k get po
NAME    READY   STATUS    RESTARTS   AGE
nginx   1/1     Running   0          15m
```

but not for kube-system:

```
KUBECONFIG=vault-kubeconfig.yml k get po -n kube-config
Error from server (Forbidden): pods is forbidden: User "system:serviceaccount:default:tmp-sa" cannot list resource "pods" in API group "" in the namespace "kube-config"
```


# teardown
```bash
kind delete cluster
kill -9 $(pgrep -x vault)
```
