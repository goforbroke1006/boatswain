package blockchain

import (
	"time"

	"github.com/goforbroke1006/boatswain/domain"
)

func NewBlockChain(storage domain.BlockStorage) domain.BlockChain {
	bc := &blockChain{storage: storage}
	bc.start()

	blocks, err := storage.LoadLastN()
	if err != nil {
		panic(err)
	}

	bc.chain = append(bc.chain, blocks...)

	return bc
}

var _ domain.BlockChain = &blockChain{}

type blockChain struct {
	chain   []*domain.Block
	storage domain.BlockStorage
}

func (bc *blockChain) Verify(block *domain.Block) error {
	//TODO implement me
	panic("implement me")
}

func (bc *blockChain) start() {
	bc.chain = append(bc.chain, genesis())
}

func (bc *blockChain) last() *domain.Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *blockChain) Append(block *domain.Block) error {
	bc.chain = append(bc.chain, block)

	err := bc.storage.Store(block)
	if err != nil {
		bc.chain = bc.chain[:len(bc.chain)-1]
		return err
	}

	return nil
}

func (bc *blockChain) GetBlocks() []*domain.Block {
	return bc.chain
}

func genesis() *domain.Block {
	unix := time.Date(2022, time.July, 2, 0, 0, 0, 0, time.UTC).Unix()
	b := domain.NewBlock(0, "", unix, []byte("Initial Block in the Chain"))
	return b
}
