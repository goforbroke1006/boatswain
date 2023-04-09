package blockchain

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/goforbroke1006/boatswain/domain"
)

func NewSyncer(
	storage domain.BlockStorage,
	reconOut chan<- *domain.ReconciliationReq,
	reconIn <-chan *domain.ReconciliationResp,
) *Syncer {
	return &Syncer{
		storage:     storage,
		reconOut:    reconOut,
		reconIn:     reconIn,
		blocksCount: 0,
	}
}

type Syncer struct {
	storage domain.BlockStorage

	reconOut chan<- *domain.ReconciliationReq
	reconIn  <-chan *domain.ReconciliationResp

	blocksCount uint64
}

func (s *Syncer) Run(ctx context.Context) error {
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

		zap.L().Info("reconciliation request", zap.Uint64("after", uint64(lastBlock.ID)))
		s.reconOut <- &domain.ReconciliationReq{AfterIndex: lastBlock.ID}

		for {
			var payload *domain.ReconciliationResp
			select {
			case <-ctx.Done():
				return ctx.Err()
			case payload = <-s.reconIn:
			// ok
			case <-time.After(6 * time.Second):
				zap.L().Warn("reconciliation timeout", zap.Uint64("after", uint64(lastBlock.ID)))
				return nil
			}

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

			s.blocksCount += uint64(len(payload.NextBlocks))

			break
		}
	}
}

func (s *Syncer) Count() uint64 {
	return s.blocksCount
}
