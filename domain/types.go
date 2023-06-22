package domain

import (
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/peer"
)

type TxID = uuid.UUID

type Transaction struct {
	ID        TxID
	Content   []byte
	Timestamp int64
}

type BlockIndex = uint64

type BlockHash string

type Block struct {
	ID       BlockIndex
	Hash     BlockHash
	PrevHash BlockHash
	Ts       int64
	Data     []Transaction
}

type BlockVote struct {
	Block    *Block
	Sender   peer.ID
	Checksum string
}
