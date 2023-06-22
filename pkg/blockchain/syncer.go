package blockchain

import (
	"context"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"go.uber.org/zap"

	"github.com/goforbroke1006/boatswain/domain"
)

func NewSyncer(
	storage domain.BlockStorage,
	p2pPubSub *pubsub.PubSub,
	reconReqTopic string,
	reconOut chan<- *domain.ReconciliationReq,
	reconIn <-chan *domain.ReconciliationResp,
) *Syncer {
	return &Syncer{
		storage:       storage,
		p2pPubSub:     p2pPubSub,
		reconReqTopic: reconReqTopic,
		reconOut:      reconOut,
		reconIn:       reconIn,
	}
}

type Syncer struct {
	storage domain.BlockStorage

	p2pPubSub     *pubsub.PubSub
	reconReqTopic string
	reconOut      chan<- *domain.ReconciliationReq
	reconIn       <-chan *domain.ReconciliationResp
}

func (s *Syncer) Run(ctx context.Context) error {
	count, countErr := s.storage.GetCount(ctx)
	if countErr != nil {
		return countErr
	}

	if count == 0 {
		storeErr := s.storage.Store(ctx, domain.Genesis())
		if storeErr != nil {
			return storeErr
		}
	}

	for {
		peers := s.p2pPubSub.ListPeers(s.reconReqTopic)
		if len(peers) > 0 {
			zap.L().Info("peers for reconciliation found", zap.Int("count", len(peers)))
			break
		}
		time.Sleep(time.Second)
	}

	for {
		lastBlock, lastBlockErr := s.storage.GetLast(ctx)
		if lastBlockErr != nil {
			return lastBlockErr
		}

		s.reconOut <- &domain.ReconciliationReq{AfterIndex: lastBlock.ID}
		zap.L().Debug("reconciliation request", zap.Uint64("after", lastBlock.ID))

	WaitAnswerLoop:
		for {
			var payload *domain.ReconciliationResp
			select {
			case <-ctx.Done():
				return ctx.Err()
			case payload = <-s.reconIn:
			// ok
			case <-time.After(15 * time.Second):
				zap.L().Warn("reconciliation response timeout")
				break WaitAnswerLoop
			}

			if payload.AfterIndex != lastBlock.ID {
				//zap.L().Debug("skip answer",
				//	zap.Uint64("got", payload.AfterIndex), zap.Uint64("want", lastBlock.ID))
				continue // skip message for another nodes
			}

			if len(payload.NextBlocks) == 0 {
				zap.L().Info("no newest blocks")
				time.Sleep(10 * time.Second)
				break WaitAnswerLoop
			}

			// TODO: verify IDs are correct sequence

			// TODO: verify blocks hashes before store, skip invalid blocks

			storeErr := s.storage.Store(ctx, payload.NextBlocks...)
			if storeErr != nil {
				return storeErr
			}

			zap.L().Info("reconciliation progress", zap.Int("blocks", len(payload.NextBlocks)))

			break WaitAnswerLoop
		}
	}
}
