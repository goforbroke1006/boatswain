package consensus

import "github.com/goforbroke1006/boatswain/domain"

func NewProofOfStake() *ProofOfStake {
	return &ProofOfStake{
		votesMap:       make(map[domain.BlockHash]*domain.ConsensusVotePayload, 1024),
		votesCollector: make(map[domain.BlockIndex]map[domain.BlockHash]map[string]struct{}, 1024),
	}
}

var _ domain.Consensus = (*ProofOfStake)(nil)

type ProofOfStake struct {
	votesMap       map[domain.BlockHash]*domain.ConsensusVotePayload
	votesCollector map[domain.BlockIndex]map[domain.BlockHash]map[string]struct{}
}

func (p ProofOfStake) Verify(vote *domain.ConsensusVotePayload) error {
	// TODO: implement me
	return nil
}

func (p ProofOfStake) Append(vote *domain.ConsensusVotePayload, peerID string) {
	if _, hasBlock := p.votesMap[vote.Hash]; !hasBlock {
		p.votesMap[vote.Hash] = vote
	}

	if _, hasIndex := p.votesCollector[vote.Index]; !hasIndex {
		p.votesCollector[vote.Index] = make(map[domain.BlockHash]map[string]struct{})
	}
	if _, hasHash := p.votesCollector[vote.Index][vote.Hash]; !hasHash {
		p.votesCollector[vote.Index][vote.Hash] = make(map[string]struct{})
	}

	// TODO: required peer ID
	// TODO: need to modify StreamIn to return meta-data (PeerID)
	p.votesCollector[vote.Index][vote.Hash][peerID] = struct{}{}
}

func (p ProofOfStake) MakeDecision(id domain.BlockIndex) (*domain.ConsensusVotePayload, error) {
	var (
		hash       domain.BlockHash
		peersCount = 0
	)
	for optionHash := range p.votesCollector[id] {
		if len(p.votesCollector[id][optionHash]) > peersCount {
			hash = optionHash
			peersCount = len(p.votesCollector[id][optionHash])
		}
	}

	return p.votesMap[hash], nil
}

func (p ProofOfStake) Reset() {
	// TODO: clear local cache of votes
}
