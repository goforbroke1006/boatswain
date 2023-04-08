package chat

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/goforbroke1006/boatswain/domain"
	"github.com/goforbroke1006/boatswain/internal/component/dapp/chat"
)

func TestHistoryMixer(t *testing.T) {
	t.Run("negative - empty", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var (
			msgCh   = make(chan *domain.TransactionPayload)
			reconCh = make(chan *domain.ReconciliationPayload)
		)

		mixer := chat.NewHistoryMixer(10, msgCh, reconCh)
		go func() { _ = mixer.Run(ctx) }()

		// write nothing

		history := mixer.History()
		assert.Len(t, history, 0)
	})

	t.Run("positive - messages income", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var (
			msgCh   = make(chan *domain.TransactionPayload)
			reconCh = make(chan *domain.ReconciliationPayload)
		)

		mixer := chat.NewHistoryMixer(10, msgCh, reconCh)
		go func() { _ = mixer.Run(ctx) }()

		msgCh <- &domain.TransactionPayload{ID: uuid.New(), Timestamp: 1680979113}
		msgCh <- &domain.TransactionPayload{ID: uuid.New(), Timestamp: 1680979112}
		msgCh <- &domain.TransactionPayload{ID: uuid.New(), Timestamp: 1680979111}
		msgCh <- &domain.TransactionPayload{ID: uuid.New(), Timestamp: 1680979110}
		msgCh <- &domain.TransactionPayload{ID: uuid.New(), Timestamp: 1680979109}
		msgCh <- &domain.TransactionPayload{ID: uuid.New(), Timestamp: 1680979108}
		msgCh <- &domain.TransactionPayload{ID: uuid.New(), Timestamp: 1680979107}
		msgCh <- &domain.TransactionPayload{ID: uuid.New(), Timestamp: 1680979106}
		msgCh <- &domain.TransactionPayload{ID: uuid.New(), Timestamp: 1680979105}
		msgCh <- &domain.TransactionPayload{ID: uuid.New(), Timestamp: 1680979104}
		msgCh <- &domain.TransactionPayload{ID: uuid.New(), Timestamp: 1680979103}
		msgCh <- &domain.TransactionPayload{ID: uuid.New(), Timestamp: 1680979102}
		msgCh <- &domain.TransactionPayload{ID: uuid.New(), Timestamp: 1680979101}

		history := mixer.History()
		assert.Len(t, history, 10)
		assert.Equal(t, history[0].Timestamp, int64(1680979104))
		assert.Equal(t, history[9].Timestamp, int64(1680979113))
	})

	t.Run("positive - reconciliation income", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var (
			msgCh   = make(chan *domain.TransactionPayload)
			reconCh = make(chan *domain.ReconciliationPayload)
		)

		mixer := chat.NewHistoryMixer(10, msgCh, reconCh)
		go func() { _ = mixer.Run(ctx) }()

		reconCh <- &domain.ReconciliationPayload{}
		reconCh <- &domain.ReconciliationPayload{}
		reconCh <- &domain.ReconciliationPayload{}

		history := mixer.History()
		assert.Len(t, history, 10)
		assert.Equal(t, history[0].Timestamp, int64(1680979104))
		assert.Equal(t, history[9].Timestamp, int64(1680979113))
	})
}
