package consensus

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goforbroke1006/boatswain/domain"
	"github.com/goforbroke1006/boatswain/internal"
	"github.com/goforbroke1006/boatswain/pkg/consensus"
)

func TestMakeDecision(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		pos := consensus.NewProofOfStake()

		var (
			currentID    = domain.BlockIndex(123)
			previousHash = internal.GetSHA256(fmt.Sprintf("%d%s%d%s",
				122, "fake", 1234567890, "no content"))
			blockA = domain.NewBlock(currentID, domain.BlockHash(previousHash), 1680818130, []byte("content A"))
			blockB = domain.NewBlock(currentID, domain.BlockHash(previousHash), 1680818130, []byte("content B"))
		)

		pos.Append(blockA, "peer-1")
		pos.Append(blockB, "peer-2")
		pos.Append(blockA, "peer-3")
		pos.Append(blockB, "peer-4")
		pos.Append(blockA, "peer-5")

		nextBlock, err := pos.MakeDecision(currentID)
		assert.Nil(t, err)
		assert.NotNil(t, nextBlock)
		assert.Equal(t, blockA, nextBlock)
	})

	t.Run("one peer can't spam", func(t *testing.T) {
		pos := consensus.NewProofOfStake()

		var (
			currentID    = domain.BlockIndex(123)
			previousHash = internal.GetSHA256(fmt.Sprintf("%d%s%d%s",
				122, "fake", 1234567890, "no content"))
			blockA = domain.NewBlock(currentID, domain.BlockHash(previousHash), 1680818130, []byte("content A"))
			blockB = domain.NewBlock(currentID, domain.BlockHash(previousHash), 1680818130, []byte("content B"))
		)

		pos.Append(blockA, "peer-1")
		pos.Append(blockA, "peer-1")
		pos.Append(blockA, "peer-1")
		pos.Append(blockA, "peer-1")
		pos.Append(blockA, "peer-1")
		pos.Append(blockA, "peer-1")

		pos.Append(blockB, "peer-2")
		pos.Append(blockB, "peer-3")
		pos.Append(blockB, "peer-3")
		pos.Append(blockB, "peer-3")

		nextBlock, err := pos.MakeDecision(currentID)
		assert.Nil(t, err)
		assert.NotNil(t, nextBlock)
		assert.Equal(t, blockB, nextBlock)
	})
}
