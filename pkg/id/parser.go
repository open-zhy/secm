package id

import (
	"crypto"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
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

func ParsePublicKey(data []byte) (PublicKey, error) {
	// First, try to parse as DER/PEM encoded data
	block, _ := pem.Decode(data)
	if block != nil {
		data = block.Bytes
	}

	// Try PKIX format first (standard format)
	pubKey, err := x509.ParsePKIXPublicKey(data)
	if err == nil {
		switch key := pubKey.(type) {
		case *rsa.PublicKey:
			return &RSAPublicKey{key}, nil
		case *ecdh.PublicKey:
			return &ECPublicKey{key}, nil
		case *ecdsa.PublicKey:
			ecKey, err := convertECDSAPublicKeyToECDH(key)
			if err != nil {
				return nil, fmt.Errorf("failed to convert ECDSA public key to ECDH: %w", err)
			}
			return &ECPublicKey{pub: ecKey}, nil
		case ed25519.PublicKey:
			edKey, err := ecdh.X25519().NewPublicKey(key)
			if err != nil {
				return nil, fmt.Errorf("failed to create X25519 public")
			}
			return &ECPublicKey{
				pub: edKey,
			}, nil
		default:
			return nil, fmt.Errorf("unknown public key type %T", key)
		}
	}

	// If PKIX fails, try PKCS1 format for RSA keys
	if rsaKey, rsaErr := x509.ParsePKCS1PublicKey(data); rsaErr == nil {
		return &RSAPublicKey{rsaKey}, nil
	}

	// If both DER formats fail, try to parse as raw key bytes
	// Try different ECDH curves for raw bytes
	curves := []ecdh.Curve{
		ecdh.X25519(),
		ecdh.P256(),
		ecdh.P384(),
		ecdh.P521(),
	}

	for _, curve := range curves {
		if pubKey, err := curve.NewPublicKey(data); err == nil {
			return &ECPublicKey{pub: pubKey}, nil
		}
	}

	return nil, fmt.Errorf("failed to parse public key: unable to parse as DER-encoded or raw key bytes")
}

func convertECDSAPublicKeyToECDH(ecdsaPub *ecdsa.PublicKey) (*ecdh.PublicKey, error) {
	var curve ecdh.Curve

	switch ecdsaPub.Curve {
	case elliptic.P256():
		curve = ecdh.P256()
	case elliptic.P384():
		curve = ecdh.P384()
	case elliptic.P521():
		curve = ecdh.P521()
	default:
		return nil, errors.New("unsupported elliptic curve")
	}

	// Serialize ECDSA public key to uncompressed format: 0x04 || X || Y
	pubkeyBytes := elliptic.Marshal(ecdsaPub.Curve, ecdsaPub.X, ecdsaPub.Y)

	// Parse as ECDH public key
	return curve.NewPublicKey(pubkeyBytes)
}
