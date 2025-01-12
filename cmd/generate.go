package cmd

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate key",
	RunE:  runGenerate,
}

var (
	keyType string
	keySize int
)

func init() {
	rootCmd.AddCommand(generateCmd)

	// add option for key type
	// add option for key size
	generateCmd.PersistentFlags().StringVarP(&keyType, "type", "t", "rsa", "Key type")
	generateCmd.PersistentFlags().IntVar(&keySize, "size", 2048, "Key size")
}

func handleRsaKeyType() error {
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return fmt.Errorf("failed to generate RSA key: %w", err)
	}

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	if err := pem.Encode(os.Stdout, privateKeyPEM); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	return nil
}

func handleEd225519KeyType() error {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate Ed25519 key: %w", err)
	}

	privateKeyPEM := &pem.Block{
		Type:  "ED25519 PRIVATE KEY",
		Bytes: privateKey,
	}

	if err := pem.Encode(os.Stdout, privateKeyPEM); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	return nil
}

// implement p256 key type
func handleP256KeyType() error {
	// Generate P256 key pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate P256 key: %w", err)
	}

	privateKeyPEM := &pem.Block{
		Type:  "P-256 PRIVATE KEY",
		Bytes: privateKey.Y.Bytes(),
	}

	if err := pem.Encode(os.Stdout, privateKeyPEM); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	return nil
}

// implement p384 key type
func handleP384KeyType() error {
	// Generate P384 key pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate P256 key: %w", err)
	}

	privateKeyPEM := &pem.Block{
		Type:  "P-384 PRIVATE KEY",
		Bytes: privateKey.Y.Bytes(),
	}

	if err := pem.Encode(os.Stdout, privateKeyPEM); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	return nil
}

func runGenerate(cmd *cobra.Command, args []string) error {
	switch keyType {
	case "rsa":
		return handleRsaKeyType()
	case "ed25519":
		return handleEd225519KeyType()
	case "p256":
		return handleP256KeyType()
	case "p384":
		return handleP384KeyType()
	default:
		return fmt.Errorf("unsupported key type: %s", keyType)
	}
}
