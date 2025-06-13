package id

import (
	"io"
)

type Encrypter interface {
	Encrypt(plaintext []byte, key []byte) ([]byte, error)
}

type Decrypter interface {
	Decrypt(ciphertext []byte) ([]byte, error)
}

type EncodableKey interface {
	Encode(dst io.Writer) error
}

type PublicKey interface {
	EncodableKey
	Encrypter
	Bytes() []byte
}

type KeyPackageIdentity interface {
	EncodableKey
	Decrypter
	PublicKey() PublicKey
}
