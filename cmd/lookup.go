package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/FalcoSuessgott/kubectl-vault-login/pkg/jwt"
	"github.com/FalcoSuessgott/kubectl-vault-login/pkg/tableprinter"
	"github.com/FalcoSuessgott/kubectl-vault-login/pkg/tokencache"
	"github.com/caarlos0/env/v11"
	"github.com/spf13/cobra"
)

type LookupOptions struct {
	CacheDir string `env:"KUBECACHEDIR" envDefault:"~/.kube/cache/vault-login"`
	IsValid  bool   `env:"IS_VALID"`
	Format   string `env:"FORMAT"`
}

// nolint: funlen
func NeLookupCmd() *cobra.Command {
	o := &LookupOptions{}

	if err := env.ParseWithOptions(o, env.Options{Prefix: EnvVarPrefix}); err != nil {
		log.Fatal(err)
	}

	tokenCacheHelper := tokencache.New(o.CacheDir)

	cmd := &cobra.Command{
		Use:           "lookup",
		Short:         "lookup the current cached token",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := tokenCacheHelper.GetToken()
			if err != nil {
				return fmt.Errorf("failed to get token from cache: %w", err)
			}

			if o.IsValid && !jwt.IsExpired(token) {
				return errors.New("token is expired")
			}

			claims, err := jwt.ParseClaims(token)
			if err != nil {
				return fmt.Errorf("failed to parse claims: %w", err)
			}

			header := []string{"Description", "Claim", "Value"}
			var data [][]string

			type claim struct {
				description string
				value       interface{}
			}

			//nolint: forcetypeassert
			details := map[string]claim{
				"iss":           {"Issuer", claims["iss"]},
				"sub":           {"Subject", claims["sub"]},
				"aud":           {"Audience", claims["aud"]},
				"exp":           {"Expiration time", toTime(claims["exp"])},
				"nbf":           {"Not before", toTime(claims["nbf"])},
				"iat":           {"Issued at", toTime(claims["iat"])},
				"jti":           {"JWT ID", claims["jti"]},
				"namespace":     {"Kubernetes namespace", claims["kubernetes.io"].(map[string]interface{})["namespace"]},
				"servieaccount": {"Service Account Name", claims["kubernetes.io"].(map[string]interface{})["serviceaccount"].(map[string]interface{})["name"]},
				"uid":           {"Service Account UID", claims["kubernetes.io"].(map[string]interface{})["serviceaccount"].(map[string]interface{})["uid"]},
			}

			var mapKeys []string
			for k := range details {
				mapKeys = append(mapKeys, k)
			}

			sort.StringSlice(mapKeys).Sort()

			for _, k := range mapKeys {
				data = append(data, []string{details[k].description, k, fmt.Sprintf("%v", details[k].value)})
			}

			if o.Format == "json" {
				out, err := json.Marshal(claims)
				if err != nil {
					return fmt.Errorf("failed to marshal data: %w", err)
				}

				fmt.Println(string(out))

				return nil
			}

			tableprinter.PrintTable(header, data)

			return nil
		},
	}

	cmd.Flags().StringVarP(&o.Format, "format", "f", o.Format, "Output format (VAULT_K8S_LOGIN_FORMAT)")
	cmd.Flags().StringVar(&o.CacheDir, "cache-dir", o.CacheDir, "Directory of where to cache token for reusing it until expiry (KUBECACHEDIR)")
	cmd.Flags().BoolVar(&o.IsValid, "is-valid", o.IsValid, "Exits with 0 if the current cached token is valid and not expired (VAULT_K8S_LOGIN_IS_VALID)")

	return cmd
}

// nolint: forcetypeassert
func toTime(f interface{}) time.Time {
	return time.Unix(int64(f.(float64)), 0)
}
