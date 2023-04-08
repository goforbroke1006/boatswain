package chat

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/goforbroke1006/boatswain/domain"
	"github.com/goforbroke1006/boatswain/internal/component/dapp/chat"
)

func TestHistoryMixer(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var (
			msgCh   = make(chan *domain.TransactionPayload)
			reconCh = make(chan *domain.ReconciliationPayload)
		)

		mixer := chat.NewHistoryMixer(10, msgCh, reconCh)
		go func() { _ = mixer.Run(ctx) }()

		// TODO: write

		history, historyErr := mixer.History()
		assert.Nil(t, historyErr)
		assert.NotNil(t, history)
		assert.Len(t, history, 10)
	})
}
