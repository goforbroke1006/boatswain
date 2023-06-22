package domain

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
)

var (
	ErrNoEnoughTxsForBlock = errors.New("no enough TXs for block")
)

type TransactionCache interface {
	Append(tx Transaction)

	// GetFirstN return first N transaction or ErrNoEnoughTxsForBlock if cache has less than requested.
	GetFirstN(count uint) ([]Transaction, error)
}

type TransactionSpreadInfoService interface {
	Spread(tx Transaction) error
}

type TransactionReader interface {
	Income() <-chan Transaction
}

func GetCheckSum(b *Block, id peer.ID) string {
	return "" // TODO: implement me
}

type BlockVoteSpreadInfoService interface {
	Spread(vote BlockVote) error
}

type BlockVoteReader interface {
	Income() <-chan BlockVote
}

type VoteCollector interface {
	// Append appends new votes
	Append(vote BlockVote)

	// GetMostVoted decide which block are next
	GetMostVoted(nextID BlockIndex) (*BlockVote, error)
}
