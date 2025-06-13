package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/open-zhy/secm/pkg/plugin"
	"github.com/open-zhy/secm/pkg/screen"
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
	Version:       Version,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "Optional profile name, specifiy the workspace related to the profile")
	home, _ := os.UserHomeDir()
	rootCmd.PersistentFlags().StringVarP(&pluginsDir, "plugins-dir", "r", filepath.Join(home, ".secm/plugins"), "Directory where plugins are stored")
}

func Execute() {
	// Load plugins before executing any command
	manager := plugin.NewManager(pluginsDir)
	if err := manager.LoadAll(rootCmd); err != nil {
		screen.Printf("Warning: failed to load plugins: %v\n", err)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
