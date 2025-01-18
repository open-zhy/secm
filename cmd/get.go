package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/open-zhy/secm/internal/crypto"
	"github.com/open-zhy/secm/internal/id"
	"github.com/open-zhy/secm/internal/secret"
	"github.com/open-zhy/secm/internal/workspace"
	"github.com/spf13/cobra"
)

var (
	outputFile string
	showMeta   bool
	quiet      bool
)

var getCmd = &cobra.Command{
	Use:   "get [secret-id]",
	Short: "Retrieve a secret by its ID",
	Long: `Retrieve and decrypt a secret using its ID. The secret can be output to stdout
or saved to a file using the --output flag.`,
	Args: cobra.ExactArgs(1),
	RunE: runGet,
}

func init() {
	getCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path (optional)")
	getCmd.Flags().BoolVarP(&showMeta, "meta", "m", false, "Show secret metadata")
	getCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Only output secret value")
	rootCmd.AddCommand(getCmd)
}

func runGet(cmd *cobra.Command, args []string) error {
	secretID := args[0]

	// Load workspace
	ws, err := workspace.Load(profile)
	if err != nil {
		return fmt.Errorf("failed to load workspace: %w", err)
	}

	// Load the secret
	secretPath := ws.SecretPath(secretID + ".yml")
	s, err := secret.Load(secretPath)
	if err != nil {
		return fmt.Errorf("failed to load secret: %w", err)
	}

	// Get encrypted data
	encryptedData, err := s.GetData()
	if err != nil {
		return fmt.Errorf("failed to decode secret data: %w", err)
	}

	identity, err := id.LoadKeyFile(ws.KeyPath)
	if err != nil {
		return fmt.Errorf("failed to load identity: %w", err)
	}

	// Decrypt the data
	decryptedData, err := crypto.DecryptData(identity, encryptedData)
	if err != nil {
		return fmt.Errorf("failed to decrypt secret: %w", err)
	}

	if showMeta {
		fmt.Printf("Name: %s\n", s.Name)
		if s.Description != "" {
			fmt.Printf("Description: %s\n", s.Description)
		}
		if s.Type != "" {
			fmt.Printf("Type: %s\n", s.Type)
		}
		if len(s.Tags) > 0 {
			fmt.Printf("Tags: %s\n", strings.Join(s.Tags, ", "))
		}
		fmt.Printf("Format: %s\n", s.Format)
		fmt.Printf("Created: %s\n", s.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println("\nSecret Value:")
	}

	// Handle output
	if outputFile != "" {
		if err := os.WriteFile(outputFile, decryptedData, 0600); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		if !quiet {
			fmt.Printf("Secret written to: %s\n", outputFile)
		}
	} else if !quiet {
		fmt.Println(string(decryptedData))
	} else {
		// In quiet mode, just print the value without newline
		fmt.Print(string(decryptedData))
	}

	return nil
}
