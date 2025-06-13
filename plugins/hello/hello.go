package main

import (
	"github.com/open-zhy/secm/pkg/plugin"
	"github.com/open-zhy/secm/pkg/screen"
	"github.com/spf13/cobra"
)

var helloCommand = &cobra.Command{
	Use:   "hello",
	Short: "Prints a hello message",
	Long:  `Prints a hello message.`,
	Run: func(cmd *cobra.Command, args []string) {
		screen.Println("Hello, world!")
	},
}

type Initializable struct{}

func (i *Initializable) Initialize(cmd *cobra.Command) {
	rootCmd := plugin.ResolveMainContext(cmd)

	rootCmd.AddCommand(helloCommand)
}

var Init = Initializable{}
