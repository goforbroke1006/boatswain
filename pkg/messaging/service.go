package messaging

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

func New(ps *pubsub.PubSub) Service {

	recvCh := make(chan Message)

	/*go func() {
		topic, err := ps.Join("")
		subscription, err := topic.Subscribe()
		message, err := subscription.Next(context.TODO())
		publicKey, err := message.ReceivedFrom.ExtractPublicKey()
	}()*/

	return &service{recvCh: recvCh}
}

var _ Service = (*service)(nil)

type service struct {
	recvCh chan Message
}

func (svc service) Send(from, to peer.ID, content []byte) error {
	//TODO implement me
	panic("implement me")
}

func (svc service) Recv() <-chan Message {
	//TODO implement me
	panic("implement me")
}
