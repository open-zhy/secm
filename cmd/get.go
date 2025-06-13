package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/open-zhy/secm/pkg/screen"
	"github.com/open-zhy/secm/pkg/secret"
	"github.com/open-zhy/secm/pkg/workspace"
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

	// Decrypt the data
	decryptedData, err := ws.DecryptSecret(s)
	if err != nil {
		return fmt.Errorf("failed to decrypt secret: %w", err)
	}

	if showMeta {
		screen.Printf("Name: %s\n", s.Name)
		if s.Description != "" {
			screen.Printf("Description: %s\n", s.Description)
		}
		if s.Type != "" {
			screen.Printf("Type: %s\n", s.Type)
		}
		if len(s.Tags) > 0 {
			screen.Printf("Tags: %s\n", strings.Join(s.Tags, ", "))
		}
		screen.Printf("Format: %s\n", s.Format)
		screen.Printf("Created: %s\n", s.CreatedAt.Format("2006-01-02 15:04:05"))
		screen.Println("\nSecret Value:")
	}

	// Handle output
	if outputFile != "" {
		if err := os.WriteFile(outputFile, decryptedData, 0600); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		if !quiet {
			screen.Printf("Secret written to: %s\n", outputFile)
		}
	} else if !quiet {
		screen.Println(string(decryptedData))
	} else {
		// In quiet mode, just print the value without newline
		screen.Printf(string(decryptedData))
	}

	return nil
}
