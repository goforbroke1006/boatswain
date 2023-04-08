package chat

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/goforbroke1006/boatswain/domain"
)

func NewHistoryMixer(
	limit uint,
	msgCh <-chan *domain.TransactionPayload,
	reconCh <-chan *domain.ReconciliationPayload,
) *HistoryMixer {
	hm := &HistoryMixer{
		limit:   limit,
		msgCh:   msgCh,
		reconCh: reconCh,
		cache:   make([]*domain.TransactionPayload, 0, limit),
	}
	return hm
}

type HistoryMixer struct {
	limit   uint
	msgCh   <-chan *domain.TransactionPayload
	reconCh <-chan *domain.ReconciliationPayload

	cache   []*domain.TransactionPayload
	cacheMx sync.RWMutex
}

func (hm *HistoryMixer) History() []*domain.TransactionPayload {
	return hm.cache
}

func (hm *HistoryMixer) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case recon := <-hm.reconCh:
			hm.cacheMx.Lock()

			// TODO: sort messages by timestamp ASC
			// TODO: get min and max timestamp
			// TODO: replace all between min-max
			_ = recon

			hm.cache = hm.cache[len(hm.cache)-int(hm.limit):] // leave N last messages
			hm.cacheMx.Unlock()

		case msg, msgOpen := <-hm.msgCh:
			if !msgOpen {
				return errors.New("messages channel is closed")
			}

			hm.cacheMx.Lock()

			hm.cache = append(hm.cache, msg)
			// TODO: sort messages by timestamp ASC

			hm.cache = hm.cache[len(hm.cache)-int(hm.limit):] // leave N last messages
			hm.cacheMx.Unlock()
		}
	}
}
