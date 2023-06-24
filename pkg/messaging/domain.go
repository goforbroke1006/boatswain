package messaging

import "github.com/libp2p/go-libp2p/core/peer"

type Service interface {
	Send(from, to peer.ID, content []byte) error
	Recv() <-chan Message
}

type Message struct {
	FromID      peer.ID
	ToID        peer.ID
	FromContent []byte
	ToContent   []byte
}
