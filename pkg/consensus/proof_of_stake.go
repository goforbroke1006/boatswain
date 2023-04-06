package consensus

import "github.com/goforbroke1006/boatswain/domain"

func NewProofOfStake() *ProofOfStake {
	return &ProofOfStake{
		blocksMap: make(map[domain.BlockHash]*domain.Block, 1024),
		pool:      make(map[domain.BlockIndex]map[domain.BlockHash]map[string]struct{}, 1024),
	}
}

var _ domain.Consensus = (*ProofOfStake)(nil)

type ProofOfStake struct {
	blocksMap map[domain.BlockHash]*domain.Block
	pool      map[domain.BlockIndex]map[domain.BlockHash]map[string]struct{}
}

func (p ProofOfStake) Append(block *domain.Block, peerCode string) {
	if _, hasBlock := p.blocksMap[block.Hash()]; !hasBlock {
		p.blocksMap[block.Hash()] = block
	}

	if _, hasIndex := p.pool[block.Index()]; !hasIndex {
		p.pool[block.Index()] = make(map[domain.BlockHash]map[string]struct{})
	}
	if _, hasHash := p.pool[block.Index()][block.Hash()]; !hasHash {
		p.pool[block.Index()][block.Hash()] = make(map[string]struct{})
	}

	p.pool[block.Index()][block.Hash()][peerCode] = struct{}{}
}

func (p ProofOfStake) MakeDecision(id domain.BlockIndex) (*domain.Block, error) {
	var (
		hash       domain.BlockHash
		peersCount = 0
	)
	for optionHash := range p.pool[id] {
		if len(p.pool[id][optionHash]) > peersCount {
			hash = optionHash
			peersCount = len(p.pool[id][optionHash])
		}
	}

	return p.blocksMap[hash], nil
}
