package plugin

import (
	"github.com/spf13/cobra"
)

// ResolveMainContext gets the main context of the cli command
// from any context execution of any plugin.
// It's the main func to be used to bind a plugin to the main
// command cli context
func ResolveMainContext(cmd *cobra.Command) *cobra.Command {
	rootCmd := cmd
	for {
		if rootCmd.Parent() == nil {
			break
		}
		rootCmd = rootCmd.Parent()
	}

	return rootCmd
}
