package cmd

import (
	"runtime"

	"github.com/open-zhy/secm/pkg/screen"
	"github.com/spf13/cobra"
)

var (
	Version string
	Build   string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version",
	Run: func(cmd *cobra.Command, args []string) {
		screen.Printf("secm version %s, Build=%s\n", Version, Build)
		screen.Printf("go version %s\n", runtime.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
