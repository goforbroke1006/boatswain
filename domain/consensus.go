package domain

type Consensus interface {
	// start voting

	// Append appends new votes
	Append(vote *Block, peerID string)

	// MakeDecision decide which block are next
	MakeDecision(id BlockIndex) (*Block, error)
}
