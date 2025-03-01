package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ParseExpiry takes a jwt and returns the exp claim.
// nolint: perfsprint
func ParseExpiry(token string) (time.Time, error) {
	t, _, err := jwt.NewParser().ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing jwt: %w", err)
	}

	// Extract exp claim
	if claims, ok := t.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			return time.Unix(int64(exp), 0), nil
		}

		return time.Time{}, fmt.Errorf("exp claim not found or invalid")
	}

	return time.Time{}, fmt.Errorf("invalid claim format")
}

func ParseClaims(token string) (map[string]interface{}, error) {
	t, _, err := jwt.NewParser().ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("error parsing jwt: %w", err)
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, errors.New("error parsing claims")
}

func IsExpired(token string) bool {
	if exp, err := ParseExpiry(token); err == nil {
		return time.Now().Before(exp)
	}

	return false
}
