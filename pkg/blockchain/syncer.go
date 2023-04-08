package blockchain

import (
	"context"
	"time"

	"github.com/goforbroke1006/boatswain/domain"
)

func NewSyncer(bc domain.BlockChain, storage domain.BlockStorage) *Syncer {
	return &Syncer{
		bc:      bc,
		storage: storage,
	}
}

type Syncer struct {
	bc      domain.BlockChain
	storage domain.BlockStorage
}

func (s Syncer) Init(ctx context.Context) error {
	blocks, err := s.storage.LoadLastN(1024)
	if err != nil {
		return err
	}

	s.bc.Append(blocks)
}

func (s Syncer) Run(ctx context.Context) error {
	timer := time.NewTimer(time.Minute)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			// TODO: flush blocks to DB
		}
	}
}
