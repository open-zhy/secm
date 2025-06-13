package transfer

import (
	"context"
	"fmt"
	"time"

	"github.com/open-zhy/secm/pkg/workspace"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
)

type NodeOption struct {
	Timeout    time.Duration
	Port       int
	Context    context.Context
	CancelFunc context.CancelFunc
}

type Node interface {
	host.Host
	Close() error
}

type TransfererNode interface {
	Node
	HandleSecretTransfer(s network.Stream, ws *workspace.Workspace, payload *SecretTransferPayload)
}

type ReceiverNode interface {
	Node
	HandleSecretReceive(s network.Stream, ws *workspace.Workspace)
	PublicKey() []byte
}

type PeerId struct {
	crypto.PrivKey
	crypto.PubKey
}

func createLocalNode(opt *NodeOption) (*TransferStreamerNode, error) {
	// Generate a key pair for this host. We will use it at least
	// to obtain a valid host ID.
	priv, pub, err := crypto.GenerateKeyPair(
		crypto.Ed25519, // Select your key type. Ed25519 are nice short
		-1,             // Select key length when possible (i.e. RSA).
	)
	if err != nil {
		return nil, err
	}

	connmgr, err := connmgr.NewConnManager(
		1, // Lowwater
		2, // HighWater,
		connmgr.WithGracePeriod(opt.Timeout),
	)
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		// Use the keypair we generated
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", opt.Port),
		),
		// support TLS connections
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		// support noise connections
		libp2p.Security(noise.ID, noise.New),
		// support any other default transports (TCP)
		libp2p.DefaultTransports,
		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(connmgr),
		// Attempt to open ports using uPNP for NATed hosts.
		libp2p.NATPortMap(),
		// If you want to help other peers to figure out if they are behind
		// NATs, you can launch the server-side of AutoNAT too (AutoRelay
		// already runs the client)
		//
		// This service is highly rate-limited and should not cause any
		// performance issues.
		libp2p.EnableNATService(),
	}

	innerNode, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}

	// create the identity
	peerId := PeerId{
		PrivKey: priv,
		PubKey:  pub,
	}

	return &TransferStreamerNode{
		innerNode,
		opt.Context,
		opt.CancelFunc,
		peerId,
	}, nil
}
