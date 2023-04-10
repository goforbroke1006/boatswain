package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func NewBlock(id BlockIndex, prevHash BlockHash, ts time.Time, data []*Transaction) *Block {
	b := &Block{
		ID:       id,
		Hash:     "",
		PrevHash: prevHash,
		Ts:       ts.Unix(),
		Data:     data,
	}
	b.GenerateHash()
	return b
}

type Block struct {
	ID       BlockIndex     `json:"id"`
	Hash     BlockHash      `json:"hash"`
	PrevHash BlockHash      `json:"prev_hash"`
	Ts       int64          `json:"ts"`
	Data     []*Transaction `json:"data"`

	metaSenderPeerID string
}

func (b *Block) GenerateHash() {
	hashContent := fmt.Sprintf("%d-%s-%d", b.ID, b.PrevHash, b.Ts)

	for _, txp := range b.Data {
		hashContent += fmt.Sprintf("--%s-%s-%d-%s",
			txp.ID.String(), txp.PeerSender, txp.Timestamp, GetSHA256(txp.Content))
	}

	b.Hash = GetSHA256(hashContent)
}

func (b *Block) SetSender(peerID string) {
	b.metaSenderPeerID = peerID
}

func (b *Block) GetSender() string {
	return b.metaSenderPeerID
}

type BlockStorage interface {
	GetCount(ctx context.Context) (uint64, error)
	GetLast(ctx context.Context) (*Block, error)
	Store(ctx context.Context, blocks ...*Block) error
	LoadLast(count uint64) ([]*Block, error)
}

var Genesis = NewBlock(1, "", time.Unix(644996700, 0), []*Transaction{
	{
		ID:            uuid.Nil,
		PeerSender:    "",
		PeerRecipient: "",
		Content:       "GENESIS",
		Timestamp:     644996700,
	},
})
