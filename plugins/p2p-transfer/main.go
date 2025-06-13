package main

import (
	"time"

	"github.com/open-zhy/secm/plugins/p2p-transfer/transfer"
	"github.com/spf13/cobra"
)

var (
	peerAddr        string
	listenPort      int
	timeoutDuration string
	rootCmd         *cobra.Command
)

var transferCommand = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer a secret on top of a p2p protocol",
	RunE: func(cmd *cobra.Command, args []string) error {
		option := &transfer.NodeOption{
			Port: listenPort,
		}

		if timeoutDuration != "" {
			option.Timeout, _ = time.ParseDuration(timeoutDuration)
		}

		var peer *transfer.PeerOption
		if peerAddr != "" {
			peer = &transfer.PeerOption{
				Addr: peerAddr,
			}

			// no need to specify as argument for peer side
			args = []string{""}
		}

		profile, _ := rootCmd.Flags().GetString("profile")
		args = append(args, profile)

		return transfer.RunTransfer(cmd, args, option, peer)
	},
	SilenceErrors: true,
	SilenceUsage:  true,
}
