package domain

import (
	"fmt"

	"github.com/enescakir/emoji"

	"github.com/goforbroke1006/boatswain/internal"
)

type Block struct {
	Index        uint64
	Hash         string
	PreviousHash string
	Timestamp    int64
	Data         string
}

func (b *Block) GenerateHash() {
	b.Hash = internal.GetSHA256(fmt.Sprintf("%d%s%d%s",
		b.Index, b.PreviousHash, b.Timestamp, b.Data))
}

func (b Block) String() string {
	return fmt.Sprintf("%v: %d %v: %s %v: %d %v: %s",
		emoji.InputNumbers, b.Index,
		emoji.Locked, b.Hash,
		emoji.OneOClock, b.Timestamp,
		emoji.Clipboard, b.Data)
}

type BlockStorage interface {
	Store(b Block) error
	Load() ([]Block, error)
}
