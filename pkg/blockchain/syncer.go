package blockchain

import (
	"context"

	"github.com/google/uuid"

	"github.com/goforbroke1006/boatswain/domain"
)

func NewSyncer(
	storage domain.BlockStorage,
	reconOut chan<- *domain.ReconciliationReq,
	reconIn <-chan *domain.ReconciliationResp,
) *Syncer {
	return &Syncer{
		storage:  storage,
		reconOut: reconOut,
		reconIn:  reconIn,
	}
}

type Syncer struct {
	storage domain.BlockStorage

	reconOut chan<- *domain.ReconciliationReq
	reconIn  <-chan *domain.ReconciliationResp
}

func (s Syncer) Run(ctx context.Context) error {
	count, countErr := s.storage.GetCount(ctx)
	if countErr != nil {
		return countErr
	}

	if count == 0 {
		genesis := domain.NewBlock(1, "", 644996700, []*domain.TransactionPayload{
			{
				Blockchain:    "",
				ID:            uuid.Nil,
				PeerSender:    "",
				PeerRecipient: "",
				Content:       "GENESIS",
				Timestamp:     644996700,
			},
		})
		if storeErr := s.storage.Store(ctx, genesis); storeErr != nil {
			return storeErr
		}
	}

	for {
		lastBlock, lastBlockErr := s.storage.GetLast(ctx)
		if lastBlockErr != nil {
			return lastBlockErr
		}

		s.reconOut <- &domain.ReconciliationReq{AfterIndex: lastBlock.Index}

		for {
			payload := <-s.reconIn

			if payload.AfterIndex != lastBlock.Index {
				continue // skip message for another nodes
			}

			if len(payload.NextBlocks) == 0 {
				return nil
			}

			// TODO: verify blocks before store, skip invalid blocks

			if storeErr := s.storage.Store(ctx, payload.NextBlocks...); storeErr != nil {
				return storeErr
			}
		}
	}
}
