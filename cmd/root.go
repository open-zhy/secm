package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	profile string
)

var rootCmd = &cobra.Command{
	Use:   "secm",
	Short: "A secure secrets manager for local development",
	Long: `secm is a command line tool for managing secrets securely on your local machine.
It provides encryption and safe storage of sensitive information.`,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "Optional profile name, specifiy the workspace related to the profile")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
