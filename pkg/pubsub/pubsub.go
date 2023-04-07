package pubsub

import (
	"context"
	"encoding/json"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

func NewStream[T any](
	ctx context.Context,
	topicName string,
	pubSub *pubsub.PubSub,
	selfID peer.ID,
	ignoreSelf bool,
) (*Stream[T], error) {
	topic, topicErr := pubSub.Join(topicName)
	if topicErr != nil {
		return nil, topicErr
	}

	subscription, subErr := topic.Subscribe()
	if subErr != nil {
		return nil, subErr
	}

	return &Stream[T]{
		ctx:          ctx,
		pubSub:       pubSub,
		topic:        topic,
		subscription: subscription,
		selfID:       selfID,
		ignoreSelf:   ignoreSelf,
	}, nil
}

type Stream[T any] struct {
	ctx          context.Context
	pubSub       *pubsub.PubSub
	topic        *pubsub.Topic
	subscription *pubsub.Subscription
	selfID       peer.ID
	ignoreSelf   bool
	out          chan *T
}

func (s Stream[T]) Out() <-chan *T {
	return s.out
}

func (s Stream[T]) readLoop() {
	for {
		msg, err := s.subscription.Next(s.ctx)
		if err != nil {
			panic(err)
		}

		// only forward messages delivered by others
		if s.ignoreSelf && msg.ReceivedFrom == s.selfID {
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
		s.out <- obj
	}
}
