package transfer

import (
	"bufio"
	"context"
	"encoding/binary"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/open-zhy/secm/pkg/errors"
	"github.com/open-zhy/secm/pkg/screen"
	"github.com/open-zhy/secm/pkg/workspace"

	ma "github.com/multiformats/go-multiaddr"
)

type PeerOption struct {
	Addr string
}

type ReceiverPeer struct {
	node   ReceiverNode
	stream *bufio.ReadWriter
}

// NewPeer creates a new ReceiverPeer instance by establishing a connection
func NewReceiverPeer(ctx context.Context, ha ReceiverNode, ws *workspace.Workspace, peerAddr string) (*ReceiverPeer, error) {
	maddr, err := ma.NewMultiaddr(peerAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid multi-address %s", peerAddr)
	}

	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to extract peer Id %s", peerAddr)
	}

	if err := ha.Connect(ctx, *info); err != nil {
		return nil, errors.Wrapf(err, "failed to connect to peer %s", peerAddr)
	}

	// Add the destination's peer multiaddress in the peerstore.
	// This will be used during connection and stream creation by libp2p.
	ha.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	// Start a stream with the destination.
	// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
	stream, err := ha.NewStream(ctx, info.ID, TRANSFER_ENDPOINT)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot establish connection with peer %s", peerAddr)
	}

	screen.Successf("Established connection to peer %s\n", peerAddr)

	// start handshake phase
	handshakePayload, err := createHandshakePayload(ha, ws)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create handshake payload")
	}

	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(len(handshakePayload)))
	if _, err := stream.Write(lengthBytes); err != nil {
		return nil, errors.Wrapf(err, "failed to write handshake length")
	}

	if _, err := stream.Write(handshakePayload); err != nil {
		return nil, errors.Wrapf(err, "failed to write handshake to stream")
	}

	// Handle secret reception in a goroutine
	go ha.HandleSecretReceive(stream, ws)

	// Create a buffered stream so that read and writes are non-blocking.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	return &ReceiverPeer{ha, rw}, nil
}
