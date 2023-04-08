package domain

type BlockChain interface {
	Verify(block *Block) error
	Append(block *Block) error
	GetBlocks() []*Block
}
