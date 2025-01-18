package id

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
)

type ECPublicKey struct {
	pub *ecdh.PublicKey
}

func (kp *ECPublicKey) Encode(dst io.Writer) error {
	// print the public key in pem format
	pubBytes, err := x509.MarshalPKIXPublicKey(kp.pub)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}
	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}

	if err := pem.Encode(dst, pemBlock); err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}

	return nil
}

func (kp *ECPublicKey) Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	// Generate ephemeral key pair
	curve := kp.pub.Curve()
	ephemeral, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ephemeral key: %w", err)
	}

	// Perform ECDH to get shared secret
	sharedSecret, err := ephemeral.ECDH(kp.pub)
	if err != nil {
		return nil, fmt.Errorf("failed to perform ECDH: %w", err)
	}

	// Create AES cipher from shared secret
	block, err := aes.NewCipher(sharedSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Create GCM mode
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Encrypt the AES key
	ciphertext := aesgcm.Seal(nil, nonce, key, nil)

	// Marshal ephemeral public key
	ephemeralPubBytes, err := x509.MarshalPKIXPublicKey(ephemeral.PublicKey())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ephemeral public key: %w", err)
	}

	// Format: [ephemeral pubkey length (4 bytes)][ephemeral pubkey][nonce][encrypted key]
	pubLen := len(ephemeralPubBytes)
	result := make([]byte, 4+pubLen+len(nonce)+len(ciphertext))

	// Store pubkey length
	result[0] = byte(pubLen >> 24)
	result[1] = byte(pubLen >> 16)
	result[2] = byte(pubLen >> 8)
	result[3] = byte(pubLen)

	// Copy ephemeral public key
	copy(result[4:], ephemeralPubBytes)
	// Copy nonce
	copy(result[4+pubLen:], nonce)
	// Copy encrypted key
	copy(result[4+pubLen+len(nonce):], ciphertext)

	return result, nil
}

type ECDHIdKey struct {
	pk *ecdh.PrivateKey
}

func (k *ECDHIdKey) PublicKey() PublicKey {
	return &ECPublicKey{
		pub: k.pk.PublicKey(),
	}
}

func (k *ECDHIdKey) Encode(dst io.Writer) error {
	pemBytes, err := x509.MarshalPKCS8PrivateKey(k.pk)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %s", err)
	}
	pemBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pemBytes,
	}

	return pem.Encode(dst, pemBlock)
}

func (k *ECDHIdKey) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < 4 {
		return nil, fmt.Errorf("invalid ciphertext format: too short")
	}

	// Extract ephemeral public key length
	pubLen := int(uint32(ciphertext[0])<<24 | uint32(ciphertext[1])<<16 | uint32(ciphertext[2])<<8 | uint32(ciphertext[3]))
	if len(ciphertext) < 4+pubLen+12 { // 4 bytes length + pubkey + minimum nonce size
		return nil, fmt.Errorf("invalid ciphertext format: insufficient data")
	}

	// Extract ephemeral public key
	ephemeralPubBytes := ciphertext[4 : 4+pubLen]
	ephemeralPubKey, err := x509.ParsePKIXPublicKey(ephemeralPubBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ephemeral public key: %w", err)
	}

	ephemeralECDH, ok := ephemeralPubKey.(*ecdh.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid ephemeral key type")
	}

	// Perform ECDH to get shared secret
	sharedSecret, err := k.pk.ECDH(ephemeralECDH)
	if err != nil {
		return nil, fmt.Errorf("failed to perform ECDH: %w", err)
	}

	// Create AES cipher from shared secret
	block, err := aes.NewCipher(sharedSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM mode
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < 4+pubLen+nonceSize {
		return nil, fmt.Errorf("invalid ciphertext format: insufficient data for nonce")
	}

	// Extract nonce and encrypted key
	nonce := ciphertext[4+pubLen : 4+pubLen+nonceSize]
	encryptedKey := ciphertext[4+pubLen+nonceSize:]

	// Decrypt the AES key
	plaintext, err := aesgcm.Open(nil, nonce, encryptedKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt key: %w", err)
	}

	return plaintext, nil
}

func createECDHIdKey(pk *ecdh.PrivateKey) *ECDHIdKey {
	return &ECDHIdKey{pk: pk}
}
