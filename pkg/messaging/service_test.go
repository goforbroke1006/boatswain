package messaging

import (
	"context"
	"github.com/libp2p/go-libp2p/core/crypto"
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

		privateKey1, _, _ := crypto.GenerateKeyPair(crypto.RSA, 2048)
		host1, _ := libp2p.New(libp2p.Identity(privateKey1))

		privateKey2, _, _ := crypto.GenerateKeyPair(crypto.RSA, 2048)
		host2, _ := libp2p.New(libp2p.Identity(privateKey2))

		var (
			received1 = make([]string, 2)
			received2 = make([]string, 2)
		)

		go communicator(ctx, host1, privateKey1, host2, received1, "Hello", "How are u?")
		go communicator(ctx, host2, privateKey2, host1, received2, "What's up", "I'm fine!")

		<-time.After(5 * time.Second)
		cancel()

		assert.Equal(t, "Hello", received2[0])
		assert.Equal(t, "How are u?", received2[1])

		assert.Equal(t, "What's up", received1[0])
		assert.Equal(t, "I'm fine!", received1[1])
	})
}

func communicator(ctx context.Context, selfHost host.Host, privKey crypto.PrivKey, anotherHost host.Host, received []string, msgs ...string) {
	pubSub, _ := pubsub.NewGossipSub(ctx, selfHost)

	topic, joinErr := pubSub.Join("Test_basic")
	if joinErr != nil {
		panic(joinErr)
	}

	subscribe, subscrErr := topic.Subscribe()
	if subscrErr != nil {
		panic(subscrErr)
	}

	svc := New(ctx, selfHost, privKey, topic, subscribe)
	defer func() { _ = svc.Close() }()

	for _, msg := range msgs {
		_ = svc.Send(anotherHost.ID(), []byte(msg))
	}

	for receiveIdx := 0; receiveIdx < len(received); receiveIdx++ {
		if income, incomeErr := svc.GetNext(ctx); incomeErr != nil {
			if incomeErr == context.Canceled {
				break
			}
			received[receiveIdx] = incomeErr.Error()
		} else {
			received[receiveIdx] = string(income.Content)
		}
	}
}
