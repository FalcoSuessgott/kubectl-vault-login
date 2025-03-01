# Token Lookup

`kubectl-vault-login` will cache the token to `~/.kube/cache/vault-login/token` (change with `$KUBECACHEDIR`) and re-use the token until expiry.

You can check if **a token has been cachend and is still valid** by:

```bash
> kubectl vault-login lookup --is-valid; echo $?
[ERROR] token is expired.
exit status 1
1
```

## List Token Details
Once a token has been cached, you can list the details of it using the `lookup` subcommand:

```bash
> kubect vault-login lookup
      DESCRIPTION      |     CLAIM     |                     VALUE
-----------------------|---------------|-------------------------------------------------
  Audience             | aud           | [https://kubernetes.default.svc.cluster.local]
  Expiration time      | exp           | 2025-03-01 14:58:01 +1100 AEDT
  Issued at            | iat           | 2025-03-01 13:58:01 +1100 AEDT
  Issuer               | iss           | https://kubernetes.default.svc.cluster.local
  JWT ID               | jti           | 91c01363-8220-4215-865f-06947343bd03
  Kubernetes namespace | namespace     | default
  Not before           | nbf           | 2025-03-01 13:58:01 +1100 AEDT
  Service Account Name | servieaccount | tmp-sa
  Subject              | sub           | system:serviceaccount:default:tmp-sa
  Service Account UID  | uid           | 5acb748f-654c-4442-974b-fd0d8d649c7d
```

## JSON Format
Or its JSON representation for scripting purposes:

```bash
> kubectl vault-login --lookup -f=json | jq .
{
  "aud": [
    "https://kubernetes.default.svc.cluster.local"
  ],
  "exp": 1740801481,
  "iat": 1740797881,
  "iss": "https://kubernetes.default.svc.cluster.local",
  "jti": "91c01363-8220-4215-865f-06947343bd03",
  "kubernetes.io": {
    "namespace": "default",
    "serviceaccount": {
      "name": "tmp-sa",
      "uid": "5acb748f-654c-4442-974b-fd0d8d649c7d"
    }
  },
  "nbf": 1740797881,
  "sub": "system:serviceaccount:default:tmp-sa"
}
```

## Force re-creation of a token
Sometimes one does need a new authenticated token. In that case you can use `VAULT_K8S_LOGIN_FORCE_NEW`:


```bash
> VAULT_K8S_LOGIN_FORCE_NEW=true kubectl get pod
```
