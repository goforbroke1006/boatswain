package domain

type ConsensusVotePayload struct {
	metaSenderPeerID string

	Index        BlockIndex            `json:"index"`
	Hash         BlockHash             `json:"hash"`
	PreviousHash BlockHash             `json:"previous_hash"`
	Timestamp    int64                 `json:"timestamp"`
	Data         []*TransactionPayload `json:"data"`
}

func (cp *ConsensusVotePayload) SetSender(peerID string) {
	cp.metaSenderPeerID = peerID
}

func (cp *ConsensusVotePayload) GetSender() string {
	return cp.metaSenderPeerID
}

type Consensus interface {
	// start voting

	// Append appends new votes
	Append(vote *ConsensusVotePayload, peerID string)

	// MakeDecision decide which block are next
	MakeDecision(id BlockIndex) (*ConsensusVotePayload, error)
}
