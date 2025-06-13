package cmd

import (
	"fmt"
	"os"

	"github.com/open-zhy/secm/pkg/screen"
	"github.com/open-zhy/secm/pkg/secret"
	"github.com/open-zhy/secm/pkg/workspace"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [secret-id]",
	Short: "Delete a secret by its ID",
	Long:  `Delete a secret using its ID. This command will remove the secret from the workspace permanently.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDelete,
}

func init() {
	deleteCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Only output secret value")
	rootCmd.AddCommand(deleteCmd)
}

func runDelete(cmd *cobra.Command, args []string) error {
	secretID := args[0]

	// Load workspace
	ws, err := workspace.Load(profile)
	if err != nil {
		return fmt.Errorf("failed to load workspace: %w", err)
	}

	// Load the secret
	secretPath := ws.SecretPath(secretID + ".yml")
	_, err = secret.Load(secretPath)
	if err != nil {
		return fmt.Errorf("no such a secret: %w", err)
	}

	// add prompt for confirmation
	if !quiet {
		screen.RedBoldf("Are you sure you want to delete the secret '%s'?\nThis action cannot be undone. (yes/no):\n", secretID)
		var response string
		fmt.Scanln(&response)
		if response != "yes" {
			screen.Println("Deletion cancelled.")
			return nil
		}
	}

	// Delete the secret file
	if err := os.Remove(secretPath); err != nil {
		return fmt.Errorf("failed to delete secret file: %w", err)
	}

	return nil
}
