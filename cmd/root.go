package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/FalcoSuessgott/kubectl-vault-login/pkg/exec_credential"
	"github.com/FalcoSuessgott/kubectl-vault-login/pkg/jwt"
	"github.com/FalcoSuessgott/kubectl-vault-login/pkg/vault"
	"github.com/caarlos0/env/v11"
	"github.com/spf13/cobra"
)

var Version = ""

const MininmumTTL = 600.0

type Options struct {
	KubernetesSecretsMount string `env:"MOUNT"     envDefault:"kubernetes"`
	KubernetesSecretsRole  string `env:"ROLE"`
	KubernetesNamespace    string `env:"NAMESPACE"`

	ClusterRoleBinding bool   `env:"CRB"`
	TTL                string `env:"TTL"       envDefault:"1h"`
	Audiences          string `env:"AUDIENCES"`

	Version bool
}

// nolint: funlen, lll, dupword, perfsprint
func NewRootCmd() *cobra.Command {
	ctx := context.Background()

	o := &Options{}

	if err := env.ParseWithOptions(o, env.Options{Prefix: "VAULT_K8S_LOGIN_"}); err != nil {
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

			if o.KubernetesSecretsRole == "" {
				return fmt.Errorf("role is required")
			}

			d, err := time.ParseDuration(o.TTL)
			if err != nil {
				return fmt.Errorf("failed to parse ttl: %w", err)
			}

			// error if ttl is less than 10 minutes
			if d.Seconds() < MininmumTTL {
				return fmt.Errorf("ttl must be at least 10 minutes (600s) but was %2.f", d.Seconds())
			}

			// auth
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

			// fetch creds
			credentials, err := v.GetKubernetesCredentials(ctx)
			if err != nil {
				return fmt.Errorf("failed to get kubernetes credentials: %w", err)
			}

			exp, err := jwt.ParseExpiry(credentials.ServiceAccountToken)
			if err != nil {
				return fmt.Errorf("failed to parse jwt expiry: %w", err)
			}

			execCreds, err := kubeconfig.NewExecCredential(credentials.ServiceAccountToken, exp)
			if err != nil {
				return fmt.Errorf("failed to create exec credential: %w", err)
			}

			fmt.Println(string(execCreds))

			return nil
		},
	}

	cmd.Flags().StringVarP(&o.KubernetesSecretsMount, "mount", "m", o.KubernetesSecretsMount, "The Kubernetes secrets mount path (VAULT_K8S_LOGIN_MOUNT)")
	cmd.Flags().StringVarP(&o.KubernetesSecretsRole, "role", "r", o.KubernetesSecretsRole, "The name of the role to generate credentials for (VAULT_K8S_LOGIN_ROLE)")
	cmd.Flags().StringVarP(&o.KubernetesNamespace, "ns", "n", o.KubernetesNamespace, "The name of the Kubernetes namespace in which to generate the credentials (VAULT_K8S_LOGIN_NAMESPACE)")
	cmd.Flags().StringVarP(&o.TTL, "ttl", "t", o.TTL, "The TTL of the generated Kubernetes service account (VAULT_K8S_LOGIN_TTL)")
	cmd.Flags().BoolVarP(&o.ClusterRoleBinding, "crb", "c", o.ClusterRoleBinding, "If true, generate a ClusterRoleBinding to grant permissions across the whole cluster instead of within a namespace (VAULT_K8S_LOGIN_CRB)")
	cmd.Flags().StringVarP(&o.Audiences, "audiences", "a", o.Audiences, "A comma separated string containing the intended audiences of the generated Kubernetes service account (VAULT_K8S_LOGIN_AUDIENCES)")

	return cmd
}

// Execute invokes the command.
func Execute() error {
	if err := NewRootCmd().Execute(); err != nil {
		return fmt.Errorf("[ERROR] %w", err)
	}

	return nil
}
