package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/tokenhelper"
)

type Option func(*Vault)

type Vault struct {
	*api.Client

	kubernetesSecretsMount string
	kubernetesSecretsRole  string
	kubernetesNamespace    string

	clusterRoleBinding bool
	ttl                string
	audiences          string
}

// NewDefaultClient returns a new vault client wrapper.
func NewDefaultClient(opts ...Option) (*Vault, error) {
	c, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	thToken, err := tokenHelper()
	if err != nil {
		return nil, err
	}

	if thToken == "" {
		c.SetToken(thToken)
	}

	v := &Vault{Client: c}

	for _, opt := range opts {
		opt(v)
	}

	return v, nil
}

func tokenHelper() (string, error) {
	th, err := tokenhelper.NewInternalTokenHelper()
	if err != nil {
		return "", fmt.Errorf("error creating token helper: %w", err)
	}

	thToken, err := th.Get()
	if err != nil {
		return "", fmt.Errorf("error getting token from token helper: %w", err)
	}

	if thToken != "" {
		return thToken, nil
	}

	return "", nil
}

// NewClient returns a new vault client wrapper.
func NewClient(addr, token string) (*Vault, error) {
	cfg := &api.Config{
		Address: addr,
	}

	c, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	c.SetToken(token)

	return &Vault{Client: c}, nil
}

func WithTTL(ttl string) Option {
	return func(v *Vault) {
		v.ttl = ttl
	}
}

func WithKubernetesNamespace(ns string) Option {
	return func(v *Vault) {
		v.kubernetesNamespace = ns
	}
}

func WithKubernetesSecretsRole(role string) Option {
	return func(v *Vault) {
		v.kubernetesSecretsRole = role
	}
}

func WithKubernetesSecretsMount(mount string) Option {
	return func(v *Vault) {
		v.kubernetesSecretsMount = mount
	}
}

func WithClusterRoleBinding(crb bool) Option {
	return func(v *Vault) {
		v.clusterRoleBinding = crb
	}
}

func WithAudiences(aud string) Option {
	return func(v *Vault) {
		v.audiences = aud
	}
}
