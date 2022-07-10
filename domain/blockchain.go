package domain

type BlockChain interface {
	Generate(ts int64, data string) error
	GetBlocks() []Block
	Empty() bool
	Verify() error
}
