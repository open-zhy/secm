package transfer

import (
	"encoding/binary"

	"github.com/open-zhy/secm/pkg/errors"
	"github.com/open-zhy/secm/pkg/workspace"
)

var PayloadTypeHandshake uint8 = 0x01
var PayloadTypeProtocolError uint8 = 0x00

type ProtocolPayload interface {
	Encode() ([]byte, error)
	Decode(data []byte) error
}

type HandshakePayload struct {
	IdentityPubKey []byte `secm:"identityPubKey"`
	PeerPubKey     []byte `secm:"peerPubKey"`
}

func (p *HandshakePayload) Encode() ([]byte, error) {
	// for handshake, .Data should be formatted as follows:
	// - 1 byte for the type of the payload (handshake)
	// - 2 bytes for length of the public key of the receiver workspace
	// - followed by the public key bytes
	// - 2 bytes for length of peer key
	// - followed by the peer public key bytes
	data := make([]byte, 1+2+len(p.IdentityPubKey)+2+len(p.PeerPubKey))
	data[0] = PayloadTypeHandshake

	// Encode identity key
	idKeyLen := len(p.IdentityPubKey)
	binary.BigEndian.PutUint16(data[1:3], uint16(idKeyLen))
	copy(data[3:3+idKeyLen], p.IdentityPubKey[:])

	// Encode peer key
	peerKeyLen := len(p.PeerPubKey)
	binary.BigEndian.PutUint16(data[3+idKeyLen:5+idKeyLen], uint16(peerKeyLen))
	copy(data[5+idKeyLen:], p.PeerPubKey[:])

	return data, nil
}

func (p *HandshakePayload) Decode(data []byte) error {
	if len(data) < 4 {
		return errors.New("invalid handshake payload length")
	}

	if data[0] != PayloadTypeHandshake {
		return errors.New("invalid handshake payload type")
	}

	// Read the identity key length and bytes
	idKeyLen := binary.BigEndian.Uint16(data[1:3])
	if len(data) < int(3+idKeyLen) {
		return errors.New("invalid identity key length in handshake payload")
	}

	p.IdentityPubKey = make([]byte, idKeyLen)
	copy(p.IdentityPubKey[:], data[3:3+idKeyLen])

	// Read the peer key length and bytes
	peerKeyLen := binary.BigEndian.Uint16(data[3+idKeyLen : 5+idKeyLen])
	if len(data) < int(5+idKeyLen+peerKeyLen) {
		return errors.New("invalid peer key length in handshake payload")
	}

	p.PeerPubKey = make([]byte, peerKeyLen)
	copy(p.PeerPubKey[:], data[5+idKeyLen:5+idKeyLen+peerKeyLen])

	return nil
}

func createHandshakePayload(ha ReceiverNode, ws *workspace.Workspace) ([]byte, error) {
	identity, err := ws.LoadKey()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load identity key from %s", ws.KeyPath)
	}

	// Add the peer public key to the payload
	peerPubKeyBytes := ha.PublicKey()
	if peerPubKeyBytes == nil {
		return nil, errors.New("peer public key is nil")
	}

	payload := &HandshakePayload{
		IdentityPubKey: identity.PublicKey().Bytes(),
		PeerPubKey:     peerPubKeyBytes,
	}

	return payload.Encode()
}

func parseHandshakePayload(data []byte) (p *HandshakePayload, err error) {
	p = &HandshakePayload{}
	if err := p.Decode(data); err != nil {
		return nil, errors.Wrapf(err, "failed to decode handshake payload")
	}

	return
}
