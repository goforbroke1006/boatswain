package node

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"go.uber.org/zap"

	"github.com/goforbroke1006/boatswain/domain"
)

func NewApplication(
	p2pHost host.Host,
	txCache domain.TransactionCache,
	txReader domain.TransactionReader,
	blockStorage domain.BlockStorage,
	blockSpreadSvc domain.BlockVoteSpreadInfoService,
	requiredForBlockCount uint,
	blockVoteReader domain.BlockVoteReader,
	voteCollector domain.VoteCollector,
) *Application {
	return &Application{
		p2pHost:               p2pHost,
		txCache:               txCache,
		txReader:              txReader,
		blockStorage:          blockStorage,
		blockSpreadSvc:        blockSpreadSvc,
		requiredForBlockCount: requiredForBlockCount,
		blockVoteReader:       blockVoteReader,
		voteCollector:         voteCollector,
	}
}

type Application struct {
	p2pHost               host.Host
	txCache               domain.TransactionCache
	txReader              domain.TransactionReader
	blockStorage          domain.BlockStorage
	blockSpreadSvc        domain.BlockVoteSpreadInfoService
	requiredForBlockCount uint
	blockVoteReader       domain.BlockVoteReader
	voteCollector         domain.VoteCollector
}

func (app Application) Run(ctx context.Context) error {
	spreadBlockVoteTicker := time.NewTicker(10 * time.Second) // TODO: move to config
	defer spreadBlockVoteTicker.Stop()

	flushBlockTicket := time.NewTicker(15 * time.Second) // TODO: move to config
	defer flushBlockTicket.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case tx := <-app.txReader.Income():
			app.txCache.Append(tx)

		case <-spreadBlockVoteTicker.C:
			transactions, getErr := app.txCache.GetFirstN(app.requiredForBlockCount)
			if getErr != nil {
				if getErr == domain.ErrNoEnoughTxsForBlock {
					continue
				} else {
					zap.L().Error("collect txs failed", zap.Error(getErr))
					continue
				}
			}

			lastBlock, lastBlockErr := app.blockStorage.GetLast(ctx)
			if lastBlockErr != nil {
				zap.L().Error("extract last block failed", zap.Error(lastBlockErr))
				continue
			}

			block := &domain.Block{
				ID:       lastBlock.ID + 1,
				Hash:     "",
				PrevHash: lastBlock.Hash,
				Ts:       time.Now().Unix(),
				Data:     transactions,
			}

			block = domain.GenerateHash(block)

			blockVote := domain.BlockVote{
				Block:    block,
				Sender:   app.p2pHost.ID(),
				Checksum: domain.GetCheckSum(block, app.p2pHost.ID()),
			}

			if sendVoteErr := app.blockSpreadSvc.Spread(blockVote); sendVoteErr != nil {
				zap.L().Error("extract last block failed", zap.Error(lastBlockErr))
				continue
			}

		case vote := <-app.blockVoteReader.Income():
			app.voteCollector.Append(vote)

		case <-flushBlockTicket.C:
			lastBlock, lastBlockErr := app.blockStorage.GetLast(ctx)
			if lastBlockErr != nil {
				zap.L().Error("extract last block failed", zap.Error(lastBlockErr))
				continue
			}

			voted, voteErr := app.voteCollector.GetMostVoted(lastBlock.ID + 1)
			if voteErr != nil {
				zap.L().Error("extract most voted failed", zap.Error(voteErr))
				continue
			}

			if storeErr := app.blockStorage.Store(ctx, voted.Block); storeErr != nil {
				zap.L().Error("store new block failed", zap.Error(voteErr))
				continue
			}
		}
	}
}
