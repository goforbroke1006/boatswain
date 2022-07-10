package blockchain

import (
	"time"

	"github.com/goforbroke1006/boatswain/domain"
)

func NewBlockChain(storage domain.BlockStorage) domain.BlockChain {
	bc := &blockChain{storage: storage}
	bc.start()

	blocks, err := storage.Load()
	if err != nil {
		panic(err)
	}
	bc.chain = append(bc.chain, blocks...)

	return bc
}

var _ domain.BlockChain = &blockChain{}

type blockChain struct {
	chain   []domain.Block
	storage domain.BlockStorage
}

func (bc *blockChain) Empty() bool {
	return len(bc.chain) <= 1
}

func (bc *blockChain) start() {
	bc.chain = append(bc.chain, genesis())
}

func (bc *blockChain) last() domain.Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *blockChain) Generate(ts int64, data string) error {
	prevHash := bc.last().Hash
	b := domain.Block{
		Index:        uint64(len(bc.chain)),
		PreviousHash: prevHash,
		Timestamp:    ts,
		Data:         data,
	}
	b.GenerateHash()
	bc.chain = append(bc.chain, b)

	err := bc.storage.Store(b)
	if err != nil {
		bc.chain = bc.chain[:len(bc.chain)-1]
		return err
	}

	return nil
}

func (bc blockChain) GetBlocks() []domain.Block {
	return bc.chain
}

func (bc *blockChain) Verify() error {
	// TODO implement me
	return nil
}

func genesis() domain.Block {
	unix := time.Date(2022, time.July, 2, 0, 0, 0, 0, time.UTC).Unix()
	b := domain.Block{
		Index:        0,
		PreviousHash: "",
		Timestamp:    unix,
		Data:         "Initial Block in the Chain",
	}
	b.GenerateHash()
	return b
}
