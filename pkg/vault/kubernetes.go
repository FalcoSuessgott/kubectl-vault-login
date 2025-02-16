package vault

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// nolint: gosec
const kubernetesCredsAPIPath = "%s/creds/%s"

type Credentials struct {
	ServiceAccountToken     string `json:"service_account_token"`
	ServiceAccountName      string `json:"service_account_name"`
	ServiceAccountNamespace string `json:"service_account_namespace"`
	TTL                     int    `json:"ttl"`
	ValidFrom, ValidUntil time.Time
}

func (v *Vault) GetKubernetesCredentials(ctx context.Context) (*Credentials, error) {
	path := fmt.Sprintf(kubernetesCredsAPIPath, v.kubernetesSecretsMount, v.kubernetesSecretsRole)

	opts := make(map[string]interface{})

	if v.kubernetesNamespace != "" {
		opts["kubernetes_namespace"] = v.kubernetesNamespace
	}

	if v.clusterRoleBinding {
		opts["cluster_role_binding"] = v.clusterRoleBinding
	}

	if v.ttl != "" {
		opts["ttl"] = v.ttl
	}

	if v.audiences != "" {
		opts["audiences"] = v.audiences
	}

	resp, err := v.Logical().WriteWithContext(ctx, path, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get kubernetes credentials: %w", err)
	}

	if resp.Data == nil {
		return nil, errors.New("no data in response")
	}

	saToken, ok := resp.Data["service_account_token"].(string)
	if !ok {
		return nil, errors.New("no service_account_token in response")
	}

	saName, ok := resp.Data["service_account_name"].(string)
	if !ok {
		return nil, errors.New("no service_account_name in response")
	}

	saNamespace, ok := resp.Data["service_account_namespace"].(string)
	if !ok {
		return nil, errors.New("no service_account_namespace in response")
	}

	return &Credentials{
		TTL:                     resp.LeaseDuration,
		ServiceAccountName:      saName,
		ServiceAccountNamespace: saNamespace,
		ServiceAccountToken:     saToken,
	}, nil
}
