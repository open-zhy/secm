package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"

	"github.com/open-zhy/secm/pkg/id"
)

// DecryptData decrypts data that was encrypted using hybrid encryption
func DecryptData(decrypter id.Decrypter, encryptedData []byte) ([]byte, error) {
	if len(encryptedData) < 4 {
		return nil, fmt.Errorf("invalid encrypted data format")
	}

	// Extract the RSA-encrypted AES key length
	keyLen := int(binary.BigEndian.Uint32(encryptedData[:4]))
	if len(encryptedData) < 4+keyLen+12 { // 4 bytes length + key + minimum nonce size
		return nil, fmt.Errorf("invalid encrypted data format")
	}

	// Extract the RSA-encrypted AES key
	encryptedKey := encryptedData[4 : 4+keyLen]

	// Decrypt the AES key using RSA
	aesKey, err := decrypter.Decrypt(encryptedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt AES key: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM mode
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := aesgcm.NonceSize()
	if len(encryptedData) < 4+keyLen+nonceSize {
		return nil, fmt.Errorf("invalid encrypted data format")
	}

	// Extract nonce and ciphertext
	nonce := encryptedData[4+keyLen : 4+keyLen+nonceSize]
	ciphertext := encryptedData[4+keyLen+nonceSize:]

	// Decrypt the data
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return plaintext, nil
}
