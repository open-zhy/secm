package id

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
)

type RSAPublicKey struct {
	pub *rsa.PublicKey
}

func (kp *RSAPublicKey) Encode(dst io.Writer) error {
	// print the public key in pem format
	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(kp.pub),
	}

	if err := pem.Encode(dst, pemBlock); err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}

	return nil
}

func (kp *RSAPublicKey) Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, kp.pub, key)
}

type RSAIdKey struct {
	pk *rsa.PrivateKey
}

func (k *RSAIdKey) PublicKey() PublicKey {
	return &RSAPublicKey{
		pub: &k.pk.PublicKey,
	}
}

func (k *RSAIdKey) Encode(dst io.Writer) error {
	pemBytes := x509.MarshalPKCS1PrivateKey(k.pk)
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: pemBytes,
	}

	return pem.Encode(dst, pemBlock)
}

func (k *RSAIdKey) Decrypt(ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, k.pk, ciphertext)
}

func createRSAIdKey(pk *rsa.PrivateKey) *RSAIdKey {
	return &RSAIdKey{pk: pk}
}
