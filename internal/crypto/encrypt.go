package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
)

// EncryptData encrypts data using hybrid encryption (RSA + AES)
// It generates a random AES key, encrypts it with RSA, and uses it to encrypt the data
func EncryptData(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	// Generate random AES key
	aesKey := make([]byte, 32) // 256-bit key
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		return nil, fmt.Errorf("failed to generate AES key: %w", err)
	}

	// Encrypt AES key with RSA
	encryptedKey, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, aesKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt AES key: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(aesKey)
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

	// Encrypt data
	ciphertext := aesgcm.Seal(nil, nonce, data, nil)

	// Combine all parts: [RSA-encrypted AES key length (4 bytes)][RSA-encrypted AES key][nonce][ciphertext]
	keyLen := len(encryptedKey)
	result := make([]byte, 4+keyLen+len(nonce)+len(ciphertext))

	// Store key length
	result[0] = byte(keyLen >> 24)
	result[1] = byte(keyLen >> 16)
	result[2] = byte(keyLen >> 8)
	result[3] = byte(keyLen)

	// Copy encrypted key
	copy(result[4:], encryptedKey)
	// Copy nonce
	copy(result[4+keyLen:], nonce)
	// Copy ciphertext
	copy(result[4+keyLen+len(nonce):], ciphertext)

	return result, nil
}
