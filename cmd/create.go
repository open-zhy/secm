package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-zhy/secm/internal/crypto"
	"github.com/open-zhy/secm/internal/secret"
	"github.com/open-zhy/secm/internal/workspace"
	"github.com/spf13/cobra"
)

var (
	secretName   string
	secretDesc   string
	secretType   string
	secretTags   string
	secretFormat string
)

var createCmd = &cobra.Command{
	Use:   "create [file]",
	Short: "Create a new secret from a file",
	Long: `Create a new secret by encrypting the contents of a file and storing it in the secm workspace.
The file will be encrypted using the RSA identity key and stored with a unique hash identifier.`,
	Args: cobra.ExactArgs(1),
	RunE: runCreate,
}

func init() {
	createCmd.Flags().StringVarP(&secretName, "name", "n", "", "Name of the secret (required)")
	createCmd.Flags().StringVarP(&secretDesc, "description", "d", "", "Description of the secret")
	createCmd.Flags().StringVarP(&secretType, "type", "t", "", "Type of secret (e.g., api-key, certificate)")
	createCmd.Flags().StringVar(&secretTags, "tags", "", "Comma-separated list of tags")
	createCmd.Flags().StringVarP(&secretFormat, "format", "f", "text", "Format of the secret (text, json, binary)")

	createCmd.MarkFlagRequired("name")
	rootCmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// Load workspace
	ws, err := workspace.Load()
	if err != nil {
		return fmt.Errorf("failed to load workspace: %w", err)
	}

	// Read the input file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Calculate file hash
	hash := sha256.Sum256(data)
	hashStr := hex.EncodeToString(hash[:])

	// Load the private key
	privateKey, err := crypto.LoadPrivateKey(ws.KeyPath)
	if err != nil {
		return fmt.Errorf("failed to load identity key: %w", err)
	}

	// Encrypt the data using hybrid encryption
	encrypted, err := crypto.EncryptData(&privateKey.PublicKey, data)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	// Create secret with metadata
	s := secret.New(secretName, encrypted)
	s.Description = secretDesc
	s.Type = secretType
	s.Format = secretFormat
	if secretTags != "" {
		s.Tags = strings.Split(secretTags, ",")
		// Trim spaces from tags
		for i, tag := range s.Tags {
			s.Tags[i] = strings.TrimSpace(tag)
		}
	}

	// Save the secret as YAML
	secretPath := filepath.Join(ws.SecretsDir, hashStr+".yml")
	if err := s.Save(secretPath); err != nil {
		return fmt.Errorf("failed to save secret: %w", err)
	}

	fmt.Printf("Created secret '%s' with ID: %s\n", secretName, hashStr)
	fmt.Printf("Stored at: %s\n", secretPath)
	return nil
}
