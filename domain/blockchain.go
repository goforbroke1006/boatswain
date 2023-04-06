package domain

type BlockChain interface {
	Append(block *Block) error
	GetBlocks() []*Block
}
