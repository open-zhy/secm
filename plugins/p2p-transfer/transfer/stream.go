package transfer

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"path/filepath"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/open-zhy/secm/pkg/id"
	"github.com/open-zhy/secm/pkg/screen"
	"github.com/open-zhy/secm/pkg/secret"
	"github.com/open-zhy/secm/pkg/workspace"
)

// SecretTransferPayload contains the secret and its ID for transfer
type SecretTransferPayload struct {
	ID     string         `json:"secret_id"`
	Secret *secret.Secret `json:"secret"`
}

type TransferStreamerNode struct {
	host.Host
	ctx      context.Context
	teardown context.CancelFunc
	id       PeerId
}

func (n *TransferStreamerNode) PublicKey() []byte {
	if n.id.PubKey == nil {
		screen.Errorf("Public key is not set")
		return nil
	}

	pub, _ := n.id.PubKey.Raw()

	return pub
}

// NewNodeStreamerWrapper creates a new NodeStreamerWrapper instance
func (n *TransferStreamerNode) HandleSecretTransfer(s network.Stream, ws *workspace.Workspace, payload *SecretTransferPayload) {
	screen.Printf("peer connected, id=%s\n", s.ID())
	defer s.Close()
	//defer n.teardown()

	if payload == nil {
		screen.Println("No payload")
		return
	}

	lengthBytes := make([]byte, 4)
	if _, err := io.ReadFull(s, lengthBytes); err != nil {
		screen.Printf("failed to read handshake length: %s\n", err)
		return
	}

	length := binary.BigEndian.Uint32(lengthBytes)
	handshakePayload := make([]byte, length)

	if _, err := io.ReadFull(s, handshakePayload); err != nil {
		screen.Printf("failed to receive handshake: %s\n", err)
		return
	}

	// Parse the handshake payload
	p, err := parseHandshakePayload(handshakePayload)
	if err != nil {
		screen.Printf("failed to parse handshake payload: %s\n", err)
		return
	}

	// secret transfer phase
	sec := payload.Secret
	secretId := payload.ID

	if sec == nil {
		screen.Println("No secret in payload")
		return
	}

	// grant read to the receiver
	receiverPubKey, err := id.ParsePublicKey(p.IdentityPubKey)
	if err != nil {
		screen.Printf("failed to parse receiver public key: %s\n", err)
		return
	}
	sec, err = ws.Grant(receiverPubKey, sec)
	if err != nil {
		screen.Printf("failed to grant read access to receiver: %s\n", err)
		return
	}

	// For sending side - serialize payload to JSON and send
	payloadData, err := json.Marshal(payload)
	if err != nil {
		screen.Printf("Error marshaling payload: %s\n", err)
		return
	}

	if _, err := s.Write(payloadData); err != nil {
		screen.Printf("Error sending secret: %s\n", err)
		return
	}

	screen.Printf("Secret '%s' sent successfully with ID: %s\n", sec.Name, secretId)
}

// HandleSecretReceive handles receiving a secret from a peer
func (n *TransferStreamerNode) HandleSecretReceive(s network.Stream, ws *workspace.Workspace) {
	screen.Printf("Receiving secret from peer: %s\n", s.ID())
	defer s.Close()
	defer n.teardown()

	// Read all data from the stream
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, s); err != nil {
		screen.Printf("failed to receive secret data: %s\n", err)
		return
	}

	// Unmarshal the received JSON into a SecretTransferPayload struct
	var payload SecretTransferPayload
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		screen.Printf("failed to unmarshaling payload: %s\n", err)
		return
	}

	// Extract the secret and its original ID
	receiviedSecret := payload.Secret

	// Print received secret details
	screen.Printf("Received secret:\n")
	screen.Printf("  ID: ")
	screen.RedBoldf("%s\n", payload.ID)
	screen.Printf("  Name: %s\n", receiviedSecret.Name)
	screen.Printf("  Description: %s\n", receiviedSecret.Description)
	if len(receiviedSecret.Tags) > 0 {
		screen.Printf("  Tags: %v\n", receiviedSecret.Tags)
	}
	if receiviedSecret.Type != "" {
		screen.Printf("  Type: %s\n", receiviedSecret.Type)
	}
	if receiviedSecret.Format != "" {
		screen.Printf("  Format: %s\n", receiviedSecret.Format)
	}
	screen.Printf("  Created: %s\n", receiviedSecret.CreatedAt.Format("2006-01-02 15:04:05"))

	// Save the received secret to workspace with the original ID-based filename
	secretPath := filepath.Join(ws.SecretsDir, payload.ID+".yml")
	if err := receiviedSecret.Save(secretPath); err != nil {
		screen.Printf("Error saving secret to workspace: %s\n", err)
		return
	}

	screen.Printf("Secret '%s' successfully saved to workspace with ID: %s\n", receiviedSecret.Name, payload.ID)
}
