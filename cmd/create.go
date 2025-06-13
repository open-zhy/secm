package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/open-zhy/secm/pkg/crypto"
	"github.com/open-zhy/secm/pkg/errors"
	"github.com/open-zhy/secm/pkg/id"
	"github.com/open-zhy/secm/pkg/screen"
	"github.com/open-zhy/secm/pkg/secret"
	"github.com/open-zhy/secm/pkg/workspace"
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
	ws, err := workspace.Load(profile)
	if err != nil {
		return errors.Wrapf(err, "failed to load workspace")
	}

	// Read the input file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read input file")
	}

	// create uuid of the file
	secretId := uuid.NewSHA1(uuid.NameSpaceDNS, data).String()

	identity, err := id.LoadKeyFile(ws.KeyPath)
	if err != nil {
		return errors.Wrapf(err, "failed to load identity")
	}

	// Encrypt the data using hybrid encryption
	encrypted, err := crypto.EncryptData(identity.PublicKey(), data)
	if err != nil {
		return errors.Wrapf(err, "failed to encrypt data")
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
	secretPath := filepath.Join(ws.SecretsDir, secretId+".yml")
	if err := s.Save(secretPath); err != nil {
		return errors.Wrapf(err, "failed to save secret")
	}

	screen.Successf("Created secret '%s' with ID: %s\n", secretName, secretId)
	screen.Successf("Stored at: %s\n", secretPath)
	return nil
}
