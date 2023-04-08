package messaging

import (
	"context"
	"encoding/json"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// NewStreamOut creates abstraction like chan<- *T
// to marshall message to JSON,
// and write data to p2p pub-sub.
// Use generic to customize format of payload.
func NewStreamOut[T any](
	ctx context.Context,
	topicName string,
	pubSub *pubsub.PubSub,
) (*StreamOut[T], error) {
	topic, topicErr := pubSub.Join(topicName)
	if topicErr != nil {
		return nil, topicErr
	}

	s := &StreamOut[T]{
		ctx:    ctx,
		pubSub: pubSub,
		topic:  topic,

		outCh: make(chan *T, 128),
	}

	go s.writeLoop()

	return s, nil
}

type StreamOut[T any] struct {
	ctx    context.Context
	pubSub *pubsub.PubSub
	topic  *pubsub.Topic

	outCh chan *T
}

func (s StreamOut[T]) Out() chan<- *T {
	return s.outCh
}

func (s StreamOut[T]) writeLoop() {
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
