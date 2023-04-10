package cmd

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/pkg/errors"
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
		transactionTopic        = "boatswain/_transaction"
		consensusVoteTopic      = "boatswain/_vote"
		reconciliationReqTopic  = "boatswain/_reconciliation/req"
		reconciliationRespTopic = "boatswain/_reconciliation/resp"

		// discoveryServiceTag is used in our mDNS advertisements to discover other chat peers.
		discoveryServiceTag        = "github.com/goforbroke1006/boatswain/node"
		allInterfacesAnyFreePortMA = "/ip4/0.0.0.0/tcp/0"
	)

	cmd := &cobra.Command{
		Use: "node",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			// create a new libp2p Host that listens on a random TCP port
			p2pHost, p2pHostErr := libp2p.New(libp2p.ListenAddrStrings(allInterfacesAnyFreePortMA))
			if p2pHostErr != nil {
				zap.L().Fatal("p2p host listening fail", zap.Error(p2pHostErr))
			}
			zap.L().Info("host peer started", zap.String("id", p2pHost.ID().String()))

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

			txStreamIn, txStreamInErr := messaging.NewStreamIn[
				domain.Transaction,
				*domain.Transaction,
			](
				ctx, transactionTopic, p2pPubSub, p2pHost.ID(), true)
			if txStreamInErr != nil {
				zap.L().Fatal("fail", zap.Error(txStreamInErr))
			}

			voteStream, voteStreamErr := messaging.NewStreamBoth[
				domain.Block,
				*domain.Block,
			](
				ctx, consensusVoteTopic, p2pPubSub, p2pHost.ID(), true)
			if voteStreamErr != nil {
				zap.L().Fatal("fail", zap.Error(voteStreamErr))
			}

			reconStreamIn, reconStreamInErr := messaging.NewStreamIn[
				domain.ReconciliationResp,
				*domain.ReconciliationResp,
			](
				ctx, reconciliationRespTopic, p2pPubSub, p2pHost.ID(), true)
			if reconStreamInErr != nil {
				zap.L().Fatal("fail", zap.Error(reconStreamInErr))
			}
			reconStreamOut, reconStreamOutErr := messaging.NewStreamOut[domain.ReconciliationReq](
				ctx, reconciliationReqTopic, p2pPubSub)
			if reconStreamOutErr != nil {
				zap.L().Fatal("fail", zap.Error(reconStreamOutErr))
			}

			blockStorage := storage.NewBlockStorage(db)

			syncer := blockchain.NewSyncer(blockStorage, reconStreamOut.Out(), reconStreamIn.In())
			if runErr := syncer.Run(ctx); runErr != nil && !errors.Is(runErr, context.Canceled) {
				zap.L().Fatal("fail", zap.Error(runErr))
			}
			zap.L().Info("reconciliation finished", zap.Uint64("blocks", syncer.Count()))

			collector := consensus.NewNextBlockGenerator(8, txStreamIn.In(),
				blockStorage, voteStream.Out())
			go func() {
				if runErr := collector.Run(ctx); runErr != nil {
					zap.L().Fatal("fail", zap.Error(runErr))
				}
			}()

			posConsensus := consensus.NewProofOfStake()
			go func() {
				for vote := range voteStream.In() {
					if verifyErr := posConsensus.Verify(vote); verifyErr != nil {
						zap.L().Error("vote verify fail", zap.Error(verifyErr))
					}
					zap.L().Info("vote", zap.Uint64("block-id", uint64(vote.ID)))
					posConsensus.Append(vote, vote.GetSender())
				}
			}()
			go func() {
				for {
					// TODO: on cron make decision
					decision, err := posConsensus.MakeDecision(123)
					if err != nil {
						zap.L().Error("make decision fail", zap.Error(err))
						continue
					}

					if decision == nil {
						time.Sleep(10 * time.Second)
						continue
						// FIXME: not finished and produce nil-pointer panic
					}

					if storeErr := blockStorage.Store(ctx, decision); storeErr != nil {
						zap.L().Error("block store fail", zap.Error(storeErr))
						continue
					}

					posConsensus.Reset()
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
