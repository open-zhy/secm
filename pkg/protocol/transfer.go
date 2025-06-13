package protocol

import (
	"github.com/open-zhy/secm/pkg/crypto"
	"github.com/open-zhy/secm/pkg/id"
)

type TransferEnvelope struct {
	identity id.KeyPackageIdentity
}

func (t *TransferEnvelope) WrapFor(data []byte, receiver id.Encrypter) ([]byte, error) {
	// Implement the transfer wrapping logic here
	// For example, you might use the receiver's public key to encrypt the data
	return crypto.EncryptData(receiver, data)
}

func (t *TransferEnvelope) Unwrap(data []byte) ([]byte, error) {
	// Implement the transfer unwrapping logic here
	// For example, you might use the receiver's private key to decrypt the data
	return crypto.DecryptData(t.identity, data)
}

func NewTransferEnvelope(identity id.KeyPackageIdentity) *TransferEnvelope {
	return &TransferEnvelope{
		identity: identity,
	}
}
