package messaging

import (
	"context"
	"github.com/pkg/errors"

	"github.com/libp2p/go-libp2p/core/peer"
)

var (
	ErrUntrustedContentReceived = errors.New("untrusted content received")
)

type Service interface {
	Send(to peer.ID, content []byte) error
	GetNext(ctx context.Context) (Income, error)
	Close() error
}

type payloadBody struct {
	FromID      peer.ID `json:"fid,omitempty"`
	ToID        peer.ID `json:"tid,omitempty"`
	FromContent []byte  `json:"fd,omitempty"`
	ToContent   []byte  `json:"td,omitempty"`
}

type Income struct {
	Sent    peer.ID
	Content []byte
}
