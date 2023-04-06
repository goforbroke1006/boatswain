package domain

type Consensus interface {
	// start voting

	// Append appends new votes
	Append(block *Block, peerCode string)

	// MakeDecision decide which block are next
	MakeDecision(id BlockIndex) (*Block, error)
}
