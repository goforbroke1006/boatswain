package messaging

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

// NewStreamIn creates abstraction
// to read data from p2p pub-sub,
// unmarshall it, and provide as <-chan *T.
// Use generic to customize format of payload.
func NewStreamIn[Type any, PtrType interface {
	*Type
	Income
}](
	ctx context.Context,
	topicName string,
	pubSub *pubsub.PubSub,
	selfID peer.ID,
	ignoreSelf bool,
) (*StreamIn[Type, PtrType], error) {
	topic, topicErr := pubSub.Join(topicName)
	if topicErr != nil {
		return nil, topicErr
	}

	subscription, subErr := topic.Subscribe()
	if subErr != nil {
		return nil, subErr
	}

	s := &StreamIn[Type, PtrType]{
		ctx:          ctx,
		pubSub:       pubSub,
		topic:        topic,
		subscription: subscription,

		selfID:     selfID,
		ignoreSelf: ignoreSelf,

		inCh: make(chan PtrType, 128),
	}

	go s.readLoop()

	return s, nil
}

type StreamIn[Type any, PtrType interface {
	*Type
	Income
}] struct {
	ctx          context.Context
	pubSub       *pubsub.PubSub
	topic        *pubsub.Topic
	subscription *pubsub.Subscription

	selfID     peer.ID
	ignoreSelf bool

	inCh chan PtrType
}

func (s StreamIn[Type, PtrType]) In() <-chan PtrType {
	return s.inCh
}

func (s StreamIn[Type, PtrType]) readLoop() {
	for {
		msg, err := s.subscription.Next(s.ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			panic(err)
		}

		// only forward messages delivered by others
		if s.ignoreSelf && msg.ReceivedFrom == s.selfID {
			continue
		}
		var obj PtrType
		obj = new(Type)
		err = json.Unmarshal(msg.Data, obj)
		if err != nil {
			// message has invalid format
			// TODO: add warning
			continue
		}
		// send valid messages onto the Messages channel
		interface{}(obj).(Income).SetSender(msg.ReceivedFrom.String())
		s.inCh <- obj
	}
}
