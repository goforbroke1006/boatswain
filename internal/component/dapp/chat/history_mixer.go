package chat

import (
	"context"
	"sort"
	"sync"

	"github.com/pkg/errors"

	"github.com/goforbroke1006/boatswain/domain"
)

func NewHistoryMixer(
	limit uint,
	msgCh <-chan *domain.TransactionPayload,
	reconCh <-chan *domain.ReconciliationResp,
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
	reconCh <-chan *domain.ReconciliationResp

	cache   []*domain.TransactionPayload
	cacheMx sync.RWMutex
}

func (hm *HistoryMixer) History() []*domain.TransactionPayload {
	hm.cacheMx.RLock()
	defer hm.cacheMx.RUnlock()

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

			var tmp []*domain.TransactionPayload
			for _, block := range recon.NextBlocks {
				tmp = append(tmp, block.Data...)
			}
			if len(tmp) == 0 {
				break
			}
			sort.Slice(tmp, func(i, j int) bool {
				return tmp[i].Timestamp < tmp[j].Timestamp
			})

			var (
				minTs = tmp[0].Timestamp
				maxTs = tmp[len(tmp)-1].Timestamp
			)

			sort.Slice(hm.cache, func(i, j int) bool {
				return hm.cache[i].Timestamp < hm.cache[j].Timestamp
			})

			for i := len(hm.cache) - 1; i >= 0; i-- {
				if minTs <= hm.cache[i].Timestamp && hm.cache[i].Timestamp <= maxTs {
					hm.cache = append(hm.cache[:i], hm.cache[i+1:]...)
				}
			}
			hm.cache = append(hm.cache, tmp...)
			sort.Slice(hm.cache, func(i, j int) bool {
				return hm.cache[i].Timestamp < hm.cache[j].Timestamp
			})

			if len(hm.cache) > int(hm.limit) {
				hm.cache = hm.cache[len(hm.cache)-int(hm.limit):] // leave N last messages
			}
			hm.cacheMx.Unlock()

		case msg, msgOpen := <-hm.msgCh:
			if !msgOpen {
				return errors.New("messages channel is closed")
			}

			hm.cacheMx.Lock()

			hm.cache = append(hm.cache, msg)
			sort.Slice(hm.cache, func(i, j int) bool {
				return hm.cache[i].Timestamp < hm.cache[j].Timestamp
			})

			if len(hm.cache) > int(hm.limit) {
				hm.cache = hm.cache[len(hm.cache)-int(hm.limit):] // leave N last messages
			}
			hm.cacheMx.Unlock()
		}
	}
}
