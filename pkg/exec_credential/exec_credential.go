package kubeconfig

import (
	"encoding/json"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/apis/clientauthentication/v1"
)

func NewExecCredential(token string, expiry time.Time) ([]byte, error) {
	execCred := v1.ExecCredential{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "client.authentication.k8s.io/v1",
			Kind:       "ExecCredential",
		},
		Status: &v1.ExecCredentialStatus{
			Token:               token,
			ExpirationTimestamp: &metav1.Time{Time: expiry},
		},
	}

	return json.MarshalIndent(execCred, "", "\t")
}
