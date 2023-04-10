package consensus

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goforbroke1006/boatswain/domain"
	"github.com/goforbroke1006/boatswain/pkg/consensus"
)

func TestMakeDecision(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		pos := consensus.NewProofOfStake()

		var (
			ts           = int64(1680818130)
			currentID    = domain.BlockIndex(123)
			previousHash = domain.GetSHA256(fmt.Sprintf("%d%s%d%s",
				122, "fake", 1234567890, "no content"))
			voteA = &domain.Block{
				ID:       currentID,
				Hash:     domain.GetSHA256(fmt.Sprintf("%d-%s-%d-%s", currentID, previousHash, ts, "vote A")),
				PrevHash: previousHash,
				Ts:       ts,
				Data:     nil,
			}
			voteB = &domain.Block{
				ID:       currentID,
				Hash:     domain.GetSHA256(fmt.Sprintf("%d-%s-%d-%s", currentID, previousHash, ts, "vote B")),
				PrevHash: previousHash,
				Ts:       ts,
				Data:     nil,
			}
		)

		pos.Append(voteA, "peer-1")
		pos.Append(voteB, "peer-2")
		pos.Append(voteA, "peer-3")
		pos.Append(voteB, "peer-4")
		pos.Append(voteA, "peer-5")

		nextBlock, err := pos.MakeDecision(currentID)
		assert.Nil(t, err)
		assert.NotNil(t, nextBlock)
		assert.Equal(t, voteA, nextBlock)
	})

	t.Run("one peer can't spam", func(t *testing.T) {
		pos := consensus.NewProofOfStake()

		var (
			ts           = int64(1680818130)
			currentID    = domain.BlockIndex(123)
			previousHash = domain.GetSHA256(fmt.Sprintf("%d%s%d%s",
				122, "fake", 1234567890, "no content"))
			voteA = &domain.Block{
				ID:       currentID,
				Hash:     domain.GetSHA256(fmt.Sprintf("%d-%s-%d-%s", currentID, previousHash, ts, "vote A")),
				PrevHash: previousHash,
				Ts:       ts,
				Data:     nil,
			}
			voteB = &domain.Block{
				ID:       currentID,
				Hash:     domain.GetSHA256(fmt.Sprintf("%d-%s-%d-%s", currentID, previousHash, ts, "vote B")),
				PrevHash: previousHash,
				Ts:       ts,
				Data:     nil,
			}
		)

		pos.Append(voteA, "peer-1")
		pos.Append(voteA, "peer-1")
		pos.Append(voteA, "peer-1")
		pos.Append(voteA, "peer-1")
		pos.Append(voteA, "peer-1")
		pos.Append(voteA, "peer-1")

		pos.Append(voteB, "peer-2")
		pos.Append(voteB, "peer-3")
		pos.Append(voteB, "peer-3")
		pos.Append(voteB, "peer-3")

		nextBlock, err := pos.MakeDecision(currentID)
		assert.Nil(t, err)
		assert.NotNil(t, nextBlock)
		assert.Equal(t, voteB, nextBlock)
	})
}
