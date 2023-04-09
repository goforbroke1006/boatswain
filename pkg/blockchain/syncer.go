package blockchain

import (
	"context"

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
		storeErr := s.storage.Store(ctx, domain.Genesis)
		if storeErr != nil {
			return storeErr
		}
	}

	for {
		lastBlock, lastBlockErr := s.storage.GetLast(ctx)
		if lastBlockErr != nil {
			return lastBlockErr
		}

		s.reconOut <- &domain.ReconciliationReq{AfterIndex: lastBlock.ID}

		for {
			payload := <-s.reconIn

			if payload.AfterIndex != lastBlock.ID {
				continue // skip message for another nodes
			}

			if len(payload.NextBlocks) == 0 {
				return nil
			}

			// TODO: verify IDs are correct sequence

			// TODO: verify blocks hashes before store, skip invalid blocks

			storeErr := s.storage.Store(ctx, payload.NextBlocks...)
			if storeErr != nil {
				return storeErr
			}

			break
		}
	}
}
