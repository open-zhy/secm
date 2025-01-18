package crypto

import (
	"crypto/rsa"
	"fmt"

	"github.com/open-zhy/secm/internal/id"
)

// LoadPrivateKey loads the RSA private key from the given file path
func LoadPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
	pk, _, err := id.LoadPrivateKey(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load key file: %w", err)
	}

	return pk.(*rsa.PrivateKey), nil
}
