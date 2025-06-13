package transfer

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/open-zhy/secm/pkg/errors"
	"github.com/open-zhy/secm/pkg/screen"
	"github.com/open-zhy/secm/pkg/secret"
	"github.com/open-zhy/secm/pkg/workspace"
	"github.com/spf13/cobra"
)

// usage:
// the transfer command is aimed to be initiated from both
// side the sender and receiver, the only difference is set
// by `--peer` which is only specified for the peer side, which
// should be set as the address of the initiator

func RunTransfer(cmd *cobra.Command, args []string, opt *NodeOption, peerOpt *PeerOption) (err error) {
	// Load workspace
	ws, err := workspace.Load(args[1])
	if err != nil {
		return errors.Wrapf(err, "failed to load workspace")
	}

	ctx, cancel := context.WithTimeout(cmd.Context(), opt.Timeout)
	defer cancel()

	opt.Context = ctx
	opt.CancelFunc = cancel

	ha, err := createLocalNode(opt)
	if err != nil {
		return errors.Wrapf(err, "failed to initiate node")
	}

	defer ha.Close()

	screen.Infof("Node informations\n")
	screen.Infof("  Id: %s \n", ha.ID())
	screen.Infof("  Addr: %s\n", ha.Addrs())

	if peerOpt == nil {
		if err := handleInitiator(ctx, ws, ha, args[0]); err != nil {
			return errors.Wrapf(err, "failed to mount the initiator")
		}
	} else {
		if err := handleReceiverPeer(ctx, ws, ha, peerOpt.Addr); err != nil {
			return errors.Wrapf(err, "failed to connect receiver")
		}
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		screen.Infof("protocol end, closing...\n")
		return nil
	case <-sigCh:
		return nil
	}
}

func handleInitiator(_ context.Context, ws *workspace.Workspace, ha TransfererNode, secretId string) error {
	// Load the secret
	secretPath := ws.SecretPath(secretId + ".yml")
	sec, err := secret.Load(secretPath)
	if err != nil {
		return errors.Wrapf(err, "failed to load secret")
	}

	// we are in a initiator context
	screen.Printf("waiting for peer on %s...\n", TRANSFER_ENDPOINT)
	ha.SetStreamHandler(TRANSFER_ENDPOINT, func(stream network.Stream) {
		payload := &SecretTransferPayload{
			ID:     secretId,
			Secret: sec,
		}
		ha.HandleSecretTransfer(stream, ws, payload)
	})

	return nil
}

func handleReceiverPeer(ctx context.Context, ws *workspace.Workspace, ha ReceiverNode, peerAddr string) error {
	screen.Println("Connecting to peer to receive secret...")

	// initiate peer node connection - this will trigger the initiator's stream handler
	_, err := NewReceiverPeer(ctx, ha, ws, peerAddr)
	if err != nil {
		return errors.Wrapf(err, "failed to initiate peer")
	}

	screen.Println("Connected to peer, waiting for secret transfer...")

	return nil
}
