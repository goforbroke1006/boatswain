package domain

import (
	"context"
	"fmt"
)

func NewBlock(index BlockIndex, previousHash BlockHash, timestamp int64, data []*TransactionPayload) *Block {
	hashContent := fmt.Sprintf("%d-%s-%d", index, previousHash, timestamp)
	for _, txp := range data {
		hashContent += fmt.Sprintf("--%s-%s-%d-%s",
			txp.ID.String(), txp.PeerSender, txp.Timestamp, GetSHA256(txp.Content))
	}

	hash := GetSHA256(hashContent)

	return &Block{
		Index:        index,
		Hash:         hash,
		PreviousHash: previousHash,
		Timestamp:    timestamp,
		Data:         data,
	}
}

type Block struct {
	Index        BlockIndex
	Hash         BlockHash
	PreviousHash BlockHash
	Timestamp    int64
	Data         []*TransactionPayload
}

type BlockStorage interface {
	GetCount(ctx context.Context) (uint64, error)
	GetLast(ctx context.Context) (*Block, error)
	Store(ctx context.Context, b ...*Block) error
}
