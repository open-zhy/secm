package protocol

import (
	"github.com/open-zhy/secm/pkg/id"
)

type EnvelopeSender interface {
	WrapFor(data []byte, receiver id.Decrypter) ([]byte, error)
}

type EnvelopeReceiver interface {
	Unwrap(data []byte) ([]byte, error)
}
