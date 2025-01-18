package id

import (
	"crypto"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// LoadPrivateKey loads the private key from the provided file path
func LoadPrivateKey(keyPath string) (crypto.PrivateKey, string, error) {
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read key file: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, "", fmt.Errorf("failed to decode PEM block")
	}

	pk, format, err := parsePrivateKeyBytes(block.Bytes)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse private key: %w", err)
	}

	return pk, format, nil
}

// LoadKeyFile loads the identity key from the provided file path
func LoadKeyFile(keyPath string) (KeyPackageIdentity, error) {
	var identity KeyPackageIdentity
	pk, format, err := LoadPrivateKey(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	if format == "PKCS1" {
		identity = createRSAIdKey(pk.(*rsa.PrivateKey))
	} else if format == "PKCS8" {
		if dhPk, ok := pk.(*ecdh.PrivateKey); ok {
			identity = createECDHIdKey(dhPk)
		} else {
			if dsPk, ok := pk.(*ecdsa.PrivateKey); ok {
				dhPk, err := dsPk.ECDH()
				if err == nil {
					identity = createECDHIdKey(dhPk)
				} else {
					return nil, fmt.Errorf("failed to convert ECDSA key to ECDH key: %w", err)
				}
			}
		}
	}

	if identity == nil {
		return nil, fmt.Errorf("failed to resurrect private key: %s", format)
	}

	return identity, nil
}

// ParsePrivateKeyBytes attempts to parse a private key in either PKCS1 or PKCS8 format
func parsePrivateKeyBytes(der []byte) (crypto.PrivateKey, string, error) {
	// Try PKCS1 first
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, "PKCS1", nil
	}

	// If PKCS1 fails, try PKCS8
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		return key, "PKCS8", nil
	}

	return nil, "", fmt.Errorf("failed to parse private key in either PKCS1 or PKCS8 format")
}
