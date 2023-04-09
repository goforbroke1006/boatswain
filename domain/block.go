package domain

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func NewBlock(index BlockIndex, previousHash BlockHash, timestamp int64, data []*TransactionPayload) *Block {
	hashContent := fmt.Sprintf("%d-%s-%d", index, previousHash, timestamp)
	for _, txp := range data {
		hashContent += fmt.Sprintf("--%s-%s-%d-%s",
			txp.ID.String(), txp.PeerSender, txp.Timestamp, GetSHA256(txp.Content))
	}

	hash := GetSHA256(hashContent)

	return &Block{
		ID:       index,
		Hash:     hash,
		PrevHash: previousHash,
		Ts:       timestamp,
		Data:     data,
	}
}

type Block struct {
	ID       BlockIndex
	Hash     BlockHash
	PrevHash BlockHash
	Ts       int64
	Data     []*TransactionPayload
}

type BlockStorage interface {
	GetCount(ctx context.Context) (uint64, error)
	GetLast(ctx context.Context) (*Block, error)
	Store(ctx context.Context, blocks ...*Block) error
	LoadLast(count uint64) ([]*Block, error)
}

var Genesis = NewBlock(1, "", 644996700, []*TransactionPayload{
	{
		Blockchain:    "",
		ID:            uuid.Nil,
		PeerSender:    "",
		PeerRecipient: "",
		Content:       "GENESIS",
		Timestamp:     644996700,
	},
})
