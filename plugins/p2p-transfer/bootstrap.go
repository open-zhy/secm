package main

import (
	"github.com/open-zhy/secm/pkg/plugin"
	"github.com/spf13/cobra"
)

type Initializable struct{}

func (i *Initializable) Initialize(cmd *cobra.Command) {
	rootCmd = plugin.ResolveMainContext(cmd)

	transferCommand.Flags().IntVar(&listenPort, "port", 0, "Listening port of the node, 0 means random")
	transferCommand.Flags().StringVar(&peerAddr, "peer", "", "Peer address to connect with")
	transferCommand.Flags().StringVar(&timeoutDuration, "timeout", "5m", "Duration to wait incoming connection and processing the transfer")
	rootCmd.AddCommand(transferCommand)
}

var Init = Initializable{}
