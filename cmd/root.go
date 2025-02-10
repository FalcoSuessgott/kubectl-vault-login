package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/FalcoSuessgott/kubectl-vault-login/pkg/kubeconfig"
	vault "github.com/FalcoSuessgott/kubectl-vault-login/pkg/vault"
	"github.com/caarlos0/env/v11"
	"github.com/spf13/cobra"
)

const ENV_PREFIX = "VAULT_K8S_LOGIN_"

var Version = ""

type Options struct {
	KubernetesSecretsMount string `env:"KUBERNETES_SECRETS_MOUNT" envDefault:"kubernetes"`
	KubernetesSecretsRole string `env:"KUBERNETES_SECRETS_ROLE"`
	KubernetesNamespace string `env:"KUBERNETES_NAMESPACE"`

	ClusterRoleBinding bool `env:"CLUSTER_ROLE_BINDING"`
	TTL string `env:"TTL" envDefault:"1h"`
	Audiences string `env:"AUDIENCES"`

	Version bool
}

func NewRootCmd() *cobra.Command {
	ctx := context.Background()

	o := &Options{}

	if err := env.ParseWithOptions(o, env.Options{Prefix: ENV_PREFIX});  err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:           "kubectl-vault-login",
		Short:         "A kubectl plugin to to obtain access to a kubernetes cluster via HashiCorp Vaults Kubernetes secrets engine",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.Version {
				fmt.Println(Version)

				return nil
			}

			v, err := vault.NewDefaultClient(
				vault.WithKubernetesSecretsMount(o.KubernetesSecretsMount),
				vault.WithKubernetesSecretsRole(o.KubernetesSecretsRole),
				vault.WithKubernetesNamespace(o.KubernetesNamespace),
				vault.WithClusterRoleBinding(o.ClusterRoleBinding),
				vault.WithTTL(o.TTL),
				vault.WithAudiences(o.Audiences),
			)

			if err != nil {
				return fmt.Errorf("failed to create vault client: %w", err)
			}

			token, err := v.GetKubernetesCredentials(ctx)
			if err != nil {
				return fmt.Errorf("failed to get kubernetes credentials: %w", err)
			}

			execCreds, err := kubeconfig.NewExecCredential(token)
			if err != nil {
				return fmt.Errorf("failed to create exec credential: %w", err)
			}

			fmt.Println(string(execCreds))

			return nil 
		},
	}

	cmd.Flags().StringVarP(&o.KubernetesSecretsMount, "mount", "m", o.KubernetesSecretsMount, "The Kubernetes secrets mount path")
	cmd.Flags().StringVarP(&o.KubernetesSecretsRole, "role", "r", o.KubernetesSecretsRole, "The name of the role to generate credentials for")
	cmd.Flags().StringVarP(&o.KubernetesNamespace, "ns", "n", o.KubernetesNamespace, "The name of the Kubernetes namespace in which to generate the credentials")
	cmd.Flags().StringVarP(&o.TTL, "ttl", "t", o.TTL, "The ttl of the generated Kubernetes service account")
	cmd.Flags().BoolVarP(&o.ClusterRoleBinding, "crb", "c", o.ClusterRoleBinding, "If true, generate a ClusterRoleBinding to grant permissions across the whole cluster instead of within a namespace")
	cmd.Flags().StringVarP(&o.Audiences, "audiences", "a", o.Audiences, "A comma separated string containing the intended audiences of the generated Kubernetes service account")
	
	return cmd
}

// Execute invokes the command.
func Execute() error {
	if err := NewRootCmd().Execute(); err != nil {
		return fmt.Errorf("[ERROR] %w", err)
	}

	return nil
}