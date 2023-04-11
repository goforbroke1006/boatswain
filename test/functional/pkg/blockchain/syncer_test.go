package blockchain

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/goforbroke1006/boatswain/domain"
	"github.com/goforbroke1006/boatswain/pkg/blockchain"
)

func TestSyncerInit(t *testing.T) {
	t.Run("positive - with correct blocks", func(t *testing.T) {
		blockStorage := newBlockStorageSpy(t)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		t.Cleanup(cancel)

		var (
			reconReqOutCh = make(chan *domain.ReconciliationReq, 4)
			reconRestInCh = make(chan *domain.ReconciliationResp, 4)
		)

		reconRestInCh <- &domain.ReconciliationResp{
			AfterIndex: 1,
			NextBlocks: []*domain.Block{
				{ID: 2, Hash: "", PrevHash: "", Ts: 0, Data: nil},
				{ID: 3, Hash: "", PrevHash: "", Ts: 0, Data: nil},
				{ID: 4, Hash: "", PrevHash: "", Ts: 0, Data: nil},
				{ID: 5, Hash: "", PrevHash: "", Ts: 0, Data: nil},
				{ID: 6, Hash: "", PrevHash: "", Ts: 0, Data: nil},
			},
		}
		reconRestInCh <- &domain.ReconciliationResp{
			AfterIndex: 6,
			NextBlocks: []*domain.Block{
				{ID: 7, Hash: "", PrevHash: "", Ts: 0, Data: nil},
				{ID: 8, Hash: "", PrevHash: "", Ts: 0, Data: nil},
				{ID: 9, Hash: "", PrevHash: "", Ts: 0, Data: nil},
				{ID: 10, Hash: "", PrevHash: "", Ts: 0, Data: nil},
			},
		}
		reconRestInCh <- &domain.ReconciliationResp{
			AfterIndex: 10,
			NextBlocks: []*domain.Block{
				{ID: 11, Hash: "", PrevHash: "", Ts: 0, Data: nil},
				{ID: 12, Hash: "", PrevHash: "", Ts: 0, Data: nil},
			},
		}
		reconRestInCh <- &domain.ReconciliationResp{
			AfterIndex: 12,
			NextBlocks: nil,
		}

		syncer := blockchain.NewSyncer(blockStorage, reconReqOutCh, reconRestInCh)
		go func() {
			if runErr := syncer.Init(ctx); runErr != nil {
				zap.L().Fatal("fail", zap.Error(runErr))
			}
		}()

		<-time.After(time.Millisecond)

		lastBlock, _ := blockStorage.GetLast(ctx)
		assert.Equalf(t, domain.BlockIndex(12), lastBlock.ID,
			"want %d got %d", 12, lastBlock.ID)
	})

	t.Run("negative - with wrong blocks", func(t *testing.T) {
		// TODO:
	})
}

func newBlockStorageSpy(t *testing.T) *spyBlockStorage {
	return &spyBlockStorage{
		t:     t,
		cache: []*domain.Block{},
	}
}

var _ domain.BlockStorage = (*spyBlockStorage)(nil)

type spyBlockStorage struct {
	t *testing.T

	cache []*domain.Block
}

func (s *spyBlockStorage) GetCount(_ context.Context) (uint64, error) {
	count := len(s.cache)
	s.t.Logf("count %d", count)
	return uint64(count), nil
}

func (s *spyBlockStorage) GetLast(_ context.Context) (*domain.Block, error) {
	s.t.Log("last 1")
	if len(s.cache) == 0 {
		return nil, errors.New("not found")
	}
	return s.cache[len(s.cache)-1], nil
}

func (s *spyBlockStorage) Store(_ context.Context, blocks ...*domain.Block) error {
	s.t.Logf("store %d", len(blocks))
	s.cache = append(s.cache, blocks...)
	return nil
}

func (s *spyBlockStorage) LoadLast(count uint64) ([]*domain.Block, error) {
	s.t.Logf("last %d", count)

	var result []*domain.Block
	if len(s.cache) > int(count) {
		result = s.cache[len(s.cache)-int(count):]
	} else {
		result = s.cache
	}
	return result, nil
}

func (s *spyBlockStorage) LoadAfterBlock(ctx context.Context, id domain.BlockIndex, count uint64) ([]*domain.Block, error) {
	//TODO implement me
	panic("implement me")
}
