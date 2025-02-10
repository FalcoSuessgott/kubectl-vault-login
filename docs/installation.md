# Installation
## Manual
Download the artifacts from the [Release](https://github.com/FalcoSuessgott/kubectl-vault-login/releases) section

## [`krew`](https://krew.sigs.k8s.io)
tbd.

## `brew`
```
brew install falcosuessgott/tap/kubectl-vault-login`
```

## `curl`
```
version=$(curl https://api.github.com/repos/falcosuessgott/kubectl-vault-login/releases/latest -s | jq .name -r)
curl -OL "https://github.com/FalcoSuessgott/kubectl-vault-login/releases/download/${version}/vkv_$(uname)_$(uname -m).tar.gz"
tar xzf kubectl-vault-login_$(uname)_$(uname -m).tar.gz
chmod u+x kubectl-vault-login
./kubectl-vault-login version
```

## `go`
```
git clone https://github.com/FalcoSuessgott/kubectl-vault-login.git
cd kubectl-vault-login
go build
```
