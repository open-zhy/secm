package cmd

import (
	"fmt"
	"os"

	"github.com/open-zhy/secm/internal/id"
	"github.com/open-zhy/secm/internal/workspace"
	"github.com/spf13/cobra"
)

var idCmd = &cobra.Command{
	Use:   "id",
	Short: "Retrieve identity of the current profile",
	Long: `Retrieve the identity of the current profile, based on the identity stored in 
	~/.secm/identity.key or the one provided by the --profile flag.`,
	RunE: runIdCommand,
}

func init() {
	rootCmd.AddCommand(idCmd)
}

func runIdCommand(cmd *cobra.Command, args []string) error {
	// Initialize workspace
	ws, err := workspace.Load(profile)
	if err != nil {
		return fmt.Errorf("failed to initialize workspace: %w", err)
	}

	identity, err := id.LoadKeyFile(ws.KeyPath)
	if err != nil {
		return fmt.Errorf("failed to load identity: %w", err)
	}

	return identity.PublicKey().Encode(os.Stdout)
}
