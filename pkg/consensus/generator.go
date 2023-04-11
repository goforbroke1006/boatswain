package consensus

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/goforbroke1006/boatswain/domain"
)

// NewGenerator creates abstraction to collect transaction,
// and if it enough to build new block,
// then create next block, set hast and send as vote to another peers.
func NewGenerator(
	limit uint,
	txInCh <-chan *domain.Transaction,
	blockStorage domain.BlockStorage,
	voteOutCh chan<- *domain.Block,
) *NextBlockGenerator {
	return &NextBlockGenerator{
		limit:        limit,
		txInCh:       txInCh,
		blockStorage: blockStorage,
		voteOutCh:    voteOutCh,
		cache:        make([]*domain.Transaction, 0, limit),
	}
}

type NextBlockGenerator struct {
	limit        uint
	txInCh       <-chan *domain.Transaction
	blockStorage domain.BlockStorage
	voteOutCh    chan<- *domain.Block

	cache []*domain.Transaction
}

func (g NextBlockGenerator) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case tx := <-g.txInCh:
			g.cache = append(g.cache, tx)

			if len(g.cache) >= int(g.limit) {
				lastBlock, lastBlockErr := g.blockStorage.GetLast(ctx)
				if lastBlockErr != nil {
					return errors.Wrap(lastBlockErr, "get last block fail")
				}
				nextBlock := domain.NewBlock(
					lastBlock.ID+1,
					lastBlock.Hash,
					time.Now(),
					g.cache,
				)
				g.voteOutCh <- nextBlock
			}
		}
	}
}
