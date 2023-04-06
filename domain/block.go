package domain

import (
	"fmt"

	"github.com/enescakir/emoji"

	"github.com/goforbroke1006/boatswain/internal"
)

func NewBlock(index BlockIndex, previousHash BlockHash, timestamp int64, data []byte) *Block {
	hash := internal.GetSHA256(fmt.Sprintf("%d%s%d%s",
		index, previousHash, timestamp, string(data)))

	return &Block{
		Index:        index,
		Hash:         BlockHash(hash),
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
	Data         []byte
}

func (b *Block) String() string {
	return fmt.Sprintf("%v: %d %v: %s %v: %d %v: %s",
		emoji.InputNumbers, b.Index,
		emoji.Locked, b.Hash,
		emoji.OneOClock, b.Timestamp,
		emoji.Clipboard, string(b.Data))
}

type BlockStorage interface {
	Store(b *Block) error
	Load() ([]*Block, error)
}
