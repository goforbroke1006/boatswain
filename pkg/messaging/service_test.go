package messaging

import (
	"context"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/stretchr/testify/assert"
)

func Test_basic(t *testing.T) {
	t.Run("two participants", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		host1, _ := libp2p.New()
		host2, _ := libp2p.New()

		var received1 []string = make([]string, 0, 100)
		var received2 []string = make([]string, 0, 100)

		go communicator(ctx, host1, host2, received1, "Hello", "How are u?")
		go communicator(ctx, host2, host1, received2, "What's up", "I'm fine!")

		<-time.After(5 * time.Second)
		cancel()

		assert.Equal(t, "Hello", received2[0])
		assert.Equal(t, "How are u?", received2[1])

		assert.Equal(t, "What's up", received1[0])
		assert.Equal(t, "I'm fine!", received1[1])
	})
}

func communicator(ctx context.Context, selfHost, anotherHost host.Host, received []string, msgs ...string) {
	pubSub, _ := pubsub.NewGossipSub(ctx, selfHost)

	svc := New(pubSub)

	nextOutIdx := 0

	for {
		select {
		case <-ctx.Done():
			return
		case income := <-svc.Recv():
			_ = income
		case <-time.After(time.Second):
			if nextOutIdx < len(msgs) {
				msg := msgs[nextOutIdx]
				_ = svc.Send(selfHost.ID(), anotherHost.ID(), []byte(msg))
				nextOutIdx++
			}
		}
	}
}
