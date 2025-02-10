package vault

import (
	"context"
	"errors"
	"fmt"
)

// nolint: gosec
const kubernetesCredsAPIPath = "%s/creds/%s"

func (v *Vault) GetKubernetesCredentials(ctx context.Context) (string, error) {
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
		return "", fmt.Errorf("failed to get kubernetes credentials: %w", err)
	}

	if resp.Data == nil {
		return "", errors.New("no data in response")
	}

	token, ok := resp.Data["service_account_token"].(string)
	if !ok {
		return "", errors.New("no service_account_token in response")
	}

	return token, nil
}
