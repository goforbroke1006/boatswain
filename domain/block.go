package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func GenerateHash(b *Block) *Block {
	hashContent := fmt.Sprintf("%d-%s-%d", b.ID, b.PrevHash, b.Ts)

	for _, tx := range b.Data {
		hashContent += fmt.Sprintf("--%s-%d-%s",
			tx.ID.String(), tx.Timestamp, GetSHA256(tx.Content))
	}

	b.Hash = BlockHash(GetSHA256([]byte(hashContent)))

	return b
}

func Genesis() *Block {
	const (
		genesisTs        = 644996700
		genesisTxContent = "GENESIS"
	)

	block := &Block{
		ID:       1,
		Hash:     "",
		PrevHash: "",
		Ts:       genesisTs,
		Data: []Transaction{
			{
				ID:        uuid.Nil,
				Content:   []byte(genesisTxContent),
				Timestamp: genesisTs,
			},
		},
	}

	block = GenerateHash(block)

	return block
}

type BlockStorage interface {
	GetCount(ctx context.Context) (uint64, error)
	GetLast(ctx context.Context) (*Block, error)
	Store(ctx context.Context, blocks ...*Block) error
	LoadLast(count uint64) ([]*Block, error)
	LoadAfterBlock(ctx context.Context, id BlockIndex, count uint64) ([]*Block, error)
}
