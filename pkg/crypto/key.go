package crypto

import (
	"crypto/rsa"

	"github.com/open-zhy/secm/pkg/errors"
	"github.com/open-zhy/secm/pkg/id"
)

// LoadPrivateKey loads the RSA private key from the given file path
func LoadPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
	pk, _, err := id.LoadPrivateKey(keyPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load key file")
	}

	return pk.(*rsa.PrivateKey), nil
}
