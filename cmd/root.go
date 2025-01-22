package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/open-zhy/secm/internal/plugin"
	"github.com/spf13/cobra"
)

var (
	profile    string
	pluginsDir string
)

var rootCmd = &cobra.Command{
	Use:   "secm",
	Short: "A secure secrets manager for local development",
	Long: `secm is a command line tool for managing secrets securely on your local machine.
It provides encryption and safe storage of sensitive information.`,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Name() != "plugin" { // Don't load plugins during plugin management
			manager := plugin.NewManager(pluginsDir)
			if err := manager.LoadAll(); err != nil {
				fmt.Printf("Warning: failed to load plugins: %v\n", err)
			}
		}
	},
}

func init() {
	executionDir, _ := os.Getwd()
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "Optional profile name, specifiy the workspace related to the profile")
	rootCmd.PersistentFlags().StringVarP(&pluginsDir, "plugins-dir", "r", filepath.Join(executionDir, "plugins"), "Directory where plugins are stored")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
