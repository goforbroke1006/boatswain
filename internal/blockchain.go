package internal

import (
	"fmt"
	"time"

	"github.com/enescakir/emoji"
)

type block struct {
	index        uint64
	hash         string
	previousHash string
	timestamp    uint64
	data         string
}

func (b *block) generateHash() {
	b.hash = GetSHA256(fmt.Sprintf("%d%s%d%s",
		b.index, b.previousHash, b.timestamp, b.data))
}

func (b block) String() string {
	return fmt.Sprintf("%v: %d %v: %s %v: %d %v: %s",
		emoji.InputNumbers, b.index,
		emoji.Locked, b.hash,
		emoji.OneOClock, b.timestamp,
		emoji.Clipboard, b.data)
}

func genesis() block {
	unix := time.Date(2022, time.July, 2, 0, 0, 0, 0, time.UTC).Unix()
	b := block{
		index:        0,
		previousHash: "",
		timestamp:    uint64(unix),
		data:         "Initial Block in the Chain",
	}
	b.generateHash()
	return b
}

type BlockChain struct {
	chain []block
}

func (bc *BlockChain) Start() {
	bc.chain = append(bc.chain, genesis())
}

func (bc *BlockChain) last() block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *BlockChain) Generate(ts uint64, data string) {
	prevHash := bc.last().hash
	b := block{
		index:        uint64(len(bc.chain)),
		previousHash: prevHash,
		timestamp:    ts,
		data:         data,
	}
	b.generateHash()
	bc.chain = append(bc.chain, b)
}

func (bc BlockChain) GetBlocks() []block {
	return bc.chain
}
