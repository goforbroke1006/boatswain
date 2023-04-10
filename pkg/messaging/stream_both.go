package messaging

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

// NewStreamBoth creates abstraction that combine StreamIn and StreamOut
func NewStreamBoth[T Income](
	ctx context.Context,
	topicName string,
	pubSub *pubsub.PubSub,
	selfID peer.ID,
	ignoreSelf bool,
) (*StreamBoth[T], error) {
	topic, topicErr := pubSub.Join(topicName)
	if topicErr != nil {
		return nil, topicErr
	}

	subscription, subErr := topic.Subscribe()
	if subErr != nil {
		return nil, subErr
	}

	s := &StreamBoth[T]{
		ctx:          ctx,
		pubSub:       pubSub,
		topic:        topic,
		subscription: subscription,

		selfID:     selfID,
		ignoreSelf: ignoreSelf,

		inCh:  make(chan T, 128),
		outCh: make(chan T, 128),
	}

	go s.readLoop()
	go s.writeLoop()

	return s, nil
}

type StreamBoth[T Income] struct {
	ctx          context.Context
	pubSub       *pubsub.PubSub
	topic        *pubsub.Topic
	subscription *pubsub.Subscription

	selfID     peer.ID
	ignoreSelf bool

	inCh  chan T
	outCh chan T
}

func (s StreamBoth[T]) In() <-chan T {
	return s.inCh
}

func (s StreamBoth[T]) Out() chan<- T {
	return s.outCh
}

func (s StreamBoth[T]) readLoop() {
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
		var obj T
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

func (s StreamBoth[T]) writeLoop() {
WriteLoop:
	for {
		select {
		case <-s.ctx.Done():
			// TODO: print s.ctx.Err()
			break WriteLoop
		case outTx := <-s.outCh:
			payloadAsJson, _ := json.Marshal(outTx)
			if publishErr := s.topic.Publish(s.ctx, payloadAsJson); publishErr != nil {
				panic(publishErr)
			}
		}
	}

	// TODO: should I do?
	// close(s.outCh)
}
