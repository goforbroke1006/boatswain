package messaging

import (
	"context"
	"encoding/json"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

// NewStreamIn creates abstraction
// to read data from p2p pub-sub,
// unmarshall it, and provide as <-chan *T.
// Use generic to customize format of payload.
func NewStreamIn[T any](
	ctx context.Context,
	topicName string,
	pubSub *pubsub.PubSub,
	selfID peer.ID,
) (*StreamIn[T], error) {
	topic, topicErr := pubSub.Join(topicName)
	if topicErr != nil {
		return nil, topicErr
	}

	subscription, subErr := topic.Subscribe()
	if subErr != nil {
		return nil, subErr
	}

	s := &StreamIn[T]{
		ctx:          ctx,
		pubSub:       pubSub,
		topic:        topic,
		subscription: subscription,
		selfID:       selfID,

		inCh: make(chan *T, 128),
	}

	go s.readLoop()

	return s, nil
}

type StreamIn[T any] struct {
	ctx          context.Context
	pubSub       *pubsub.PubSub
	topic        *pubsub.Topic
	subscription *pubsub.Subscription
	selfID       peer.ID

	inCh chan *T
}

func (s StreamIn[T]) In() <-chan *T {
	return s.inCh
}

func (s StreamIn[T]) readLoop() {
	for {
		msg, err := s.subscription.Next(s.ctx)
		if err != nil {
			panic(err)
		}

		// only forward messages delivered by others
		if msg.ReceivedFrom == s.selfID {
			continue
		}
		obj := new(T)
		err = json.Unmarshal(msg.Data, obj)
		if err != nil {
			// message has invalid format
			// TODO: add warning
			continue
		}
		// send valid messages onto the Messages channel
		s.inCh <- obj
	}
}
