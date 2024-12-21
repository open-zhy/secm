package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/open-zhy/secm/internal/crypto"
	"github.com/open-zhy/secm/internal/workspace"
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
}

func runInit(cmd *cobra.Command, args []string) error {
	// Initialize workspace
	ws, err := workspace.Initialize()
	if err != nil {
		return fmt.Errorf("failed to initialize workspace: %w", err)
	}

	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate RSA key: %w", err)
	}

	// Save private key
	if err := crypto.SavePrivateKey(privateKey, ws.KeyPath); err != nil {
		return fmt.Errorf("failed to save identity key: %w", err)
	}

	fmt.Printf("Initialized secm workspace at %s\n", ws.RootDir)
	fmt.Printf("Generated RSA identity key at %s\n", ws.KeyPath)
	return nil
}
