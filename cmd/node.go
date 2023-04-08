package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/goforbroke1006/boatswain/domain"
	"github.com/goforbroke1006/boatswain/internal/common"
	"github.com/goforbroke1006/boatswain/internal/storage"
	"github.com/goforbroke1006/boatswain/pkg/blockchain"
	"github.com/goforbroke1006/boatswain/pkg/consensus"
	"github.com/goforbroke1006/boatswain/pkg/discovery"
	"github.com/goforbroke1006/boatswain/pkg/messaging"
)

func NewNode() *cobra.Command {
	const (
		transactionTopic    = "boatswain/transaction"
		consensusVoteTopic  = "boatswain/consensus-vote"
		reconciliationTopic = "boatswain/reconciliation"

		// discoveryServiceTag is used in our mDNS advertisements to discover other chat peers.
		discoveryServiceTag        = "github.com/goforbroke1006/boatswain/node"
		allInterfacesAnyFreePortMA = "/ip4/0.0.0.0/tcp/0"
	)

	cmd := &cobra.Command{
		Use: "node",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()

			// create a new libp2p Host that listens on a random TCP port
			p2pHost, p2pHostErr := libp2p.New(libp2p.ListenAddrStrings(allInterfacesAnyFreePortMA))
			if p2pHostErr != nil {
				zap.L().Fatal("p2p host listening fail", zap.Error(p2pHostErr))
			}

			// create a new PubSub service using the GossipSub router
			p2pPubSub, p2pPubSubErr := pubsub.NewGossipSub(ctx, p2pHost)
			if p2pPubSubErr != nil {
				zap.L().Fatal("initialize gossip sub fail", zap.Error(p2pPubSubErr))
			}

			// setup local mDNS discoverySvc
			discoverySvc := discovery.NewDiscovery(p2pHost, discoveryServiceTag)
			if discoveryErr := discoverySvc.Start(); discoveryErr != nil {
				zap.L().Fatal("initialize gossip sub fail", zap.Error(discoveryErr))
			}
			defer func() { _ = discoverySvc.Close() }()

			db, dbErr := common.OpenDBConn()
			if dbErr != nil {
				zap.L().Fatal("open db connection fail", zap.Error(dbErr))
			}
			defer func() { _ = db.Close() }()

			txStreamIn, txStreamInErr := messaging.NewStreamIn[domain.TransactionPayload](ctx, transactionTopic, p2pPubSub, p2pHost.ID())
			if txStreamInErr != nil {
				zap.L().Fatal("fail", zap.Error(txStreamInErr))
			}

			voteStreamIn, voteStreamInErr := messaging.NewStreamIn[domain.ConsensusVotePayload](ctx, consensusVoteTopic, p2pPubSub, p2pHost.ID())
			if voteStreamInErr != nil {
				zap.L().Fatal("fail", zap.Error(voteStreamInErr))
			}
			voteStreamOut, voteStreamOutErr := messaging.NewStreamOut[domain.ConsensusVotePayload](ctx, consensusVoteTopic, p2pPubSub)
			if voteStreamOutErr != nil {
				zap.L().Fatal("fail", zap.Error(voteStreamOutErr))
			}

			reconStreamOut, reconStreamOutErr := messaging.NewStreamOut[domain.ReconciliationPayload](ctx, reconciliationTopic, p2pPubSub)
			if reconStreamOutErr != nil {
				zap.L().Fatal("fail", zap.Error(reconStreamOutErr))
			}

			blockStorage := storage.NewBlockStorage(db)
			chain := blockchain.NewBlockChain(blockStorage)

			syncer := blockchain.NewSyncer(chain, blockStorage)
			if syncErr := syncer.Init(ctx); syncErr != nil {
				zap.L().Fatal("fail", zap.Error(syncErr))
			}
			go func() {
				if runErr := syncer.Run(ctx); runErr != nil {
					zap.L().Fatal("fail", zap.Error(runErr))
				}
			}()

			pos := consensus.NewProofOfStake()

			go func() {
				for tx := range txStreamIn.In() {
					// TODO: collect TX
					_ = tx

					// TODO: if TXes takes 1Mb of memory
					voteStreamOut.Out() <- &domain.ConsensusVotePayload{
						// TODO: fill with collector cache
					}
				}
			}()

			go func() {
				for vote := range voteStreamIn.In() {
					if verifyErr := pos.Verify(vote); verifyErr != nil {
						zap.L().Error("vote verify fail", zap.Error(verifyErr))
					}
					pos.Append(vote)
				}
			}()

			go func() {
				for {
					// TODO: on cron make decision
					decision, err := pos.MakeDecision(123)
					if err != nil {
						zap.L().Error("make decision fail", zap.Error(err))
						continue
					}

					block := &domain.Block{
						Index:        decision.Index,
						Hash:         decision.Hash,
						PreviousHash: decision.PreviousHash,
						Timestamp:    decision.Timestamp,
						Data:         decision.Data,
					}
					if verifyErr := chain.Append(block); verifyErr != nil {
						zap.L().Error("block verify fail", zap.Error(verifyErr))
					}
					if appendErr := chain.Append(block); appendErr != nil {
						zap.L().Error("block append fail", zap.Error(appendErr))
					}

					reconStreamOut.Out() <- &domain.ReconciliationPayload{
						// TODO: fill recon payload from block
					}
				}
			}()

			<-ctx.Done()
		},
	}

	return cmd
}

func init() {
	rootCmd.AddCommand(NewNode())
}
