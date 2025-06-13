package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/open-zhy/secm/pkg/id"
	"github.com/open-zhy/secm/pkg/screen"
	"github.com/open-zhy/secm/pkg/workspace"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize secm workspace and generate identity key",
	Long: `Initialize the secm workspace in ~/.secm directory and generate an RSA identity key
for encrypting and decrypting secrets.`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.PersistentFlags().StringVarP(&keyType, "type", "t", "rsa", "Key type, supports rsa, p256, p384, p521, ec25519")
	initCmd.PersistentFlags().IntVar(&keySize, "size", 2048, "Key size, take effect for RSA key types only")
}

func runInit(cmd *cobra.Command, args []string) error {
	// Initialize workspace
	ws, err := workspace.Initialize(profile)
	if err != nil {
		return fmt.Errorf("failed to initialize workspace: %w", err)
	}

	identity, err := id.GenerateKey(
		id.GenerateKeyOpts{
			Type: keyType,
			Size: &keySize,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	keyFile, err := os.OpenFile(ws.KeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create key file: %w", err)
	}
	defer keyFile.Close()

	if err := identity.Encode(keyFile); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	screen.Printf("Initialized secm workspace at %s\n", ws.RootDir)
	screen.Printf("Generated %s identity key at %s\n", strings.ToUpper(keyType), ws.KeyPath)
	return nil
}
