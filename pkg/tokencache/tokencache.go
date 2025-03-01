package tokencache

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/client-go/util/homedir"
)

const TokenCacheFile = "token"

// TokenCache interface for caching and obtaining tokens.
type TokenCache interface {
	// GetToken returns the token from the cache.
	GetToken() (string, error)

	// SaveToken saves the token to the cache.
	SaveToken(token []byte) error
}

// TokenCacheHelper is a helper struct for token caching.
type TokenCacheHelper struct {
	cacheDir, tokenFile string
}

// New creates a new token cache helper.
func New(cacheDir string) *TokenCacheHelper {
	if strings.HasPrefix(cacheDir, "~") {
		cacheDir = filepath.Join(homedir.HomeDir(), strings.TrimPrefix(cacheDir, "~"))
	}

	th := &TokenCacheHelper{
		cacheDir:  cacheDir,
		tokenFile: filepath.Join(cacheDir, TokenCacheFile),
	}

	return th
}

// NewTokenCache creates a new token cache.
func (th *TokenCacheHelper) GetToken() (string, error) {
	f, err := os.ReadFile(th.tokenFile)
	if err != nil {
		return "", fmt.Errorf("failed to read token file: %w", err)
	}

	return string(f), nil
}

// NewTokenCache saves a token.
func (th *TokenCacheHelper) SaveToken(token []byte) error {
	if err := os.MkdirAll(th.cacheDir, 0o700); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	if err := os.WriteFile(th.tokenFile, token, 0o600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}
