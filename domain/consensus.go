package domain

type ConsensusVotePayload struct {
	Index        BlockIndex           `json:"index"`
	Hash         BlockHash            `json:"hash"`
	PreviousHash BlockHash            `json:"previous_hash"`
	Timestamp    int64                `json:"timestamp"`
	Data         []TransactionPayload `json:"data"`
}

type Consensus interface {
	// start voting

	// Append appends new votes
	Append(vote *ConsensusVotePayload)

	// MakeDecision decide which block are next
	MakeDecision(id BlockIndex) (*ConsensusVotePayload, error)
}
