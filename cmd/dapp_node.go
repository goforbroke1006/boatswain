package cmd

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"time"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/goforbroke1006/boatswain/domain"
	"github.com/goforbroke1006/boatswain/internal/blockchain"
	"github.com/goforbroke1006/boatswain/internal/common"
	"github.com/goforbroke1006/boatswain/internal/component/util/chat"
	"github.com/goforbroke1006/boatswain/internal/storage"
	"github.com/goforbroke1006/boatswain/pkg/consensus"
)

func NewDAppNode() *cobra.Command {
	const (
		// discoveryServiceTag is used in our mDNS advertisements to discover other chat peers.
		discoveryServiceTag        = "boatswain-chat-example"
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

			// setup local mDNS discovery
			discovery := chat.NewDiscovery(p2pHost, discoveryServiceTag)
			if discoveryErr := discovery.Start(); discoveryErr != nil {
				zap.L().Fatal("initialize gossip sub fail", zap.Error(discoveryErr))
			}
			defer func() { _ = discovery.Close() }()

			db, dbErr := common.OpenDBConn()
			if dbErr != nil {
				zap.L().Fatal("open db connection fail", zap.Error(dbErr))
			}
			defer func() { _ = db.Close() }()

			pos := consensus.NewProofOfStake()

			go func() {
				topic, err := p2pPubSub.Join("consensus")
				if err != nil {
					zap.L().Fatal("join consensus topic fail", zap.Error(err))
				}

				subscription, err := topic.Subscribe()
				if err != nil {
					zap.L().Fatal("subscribe consensus topic fail", zap.Error(err))
				}

				self := p2pHost.ID()

				for {
					message, err := subscription.Next(ctx)
					if err != nil {
						zap.L().Error("read consensus message fail", zap.Error(err))
						continue
					}
					if message.ReceivedFrom == self { // only forward messages delivered by others
						continue
					}

					block := domain.Block{}
					_ = json.Unmarshal(message.Data, &block)

					pos.Append(&block, message.ReceivedFrom.String())
				}
			}()

			blockStorage := storage.NewBlockStorage(db)
			chain := blockchain.NewBlockChain(blockStorage)

			go func() {
				var nextCheckIndex domain.BlockIndex = 1
				for {
					nextBlock, err := pos.MakeDecision(nextCheckIndex)
					if err != nil {
						time.Sleep(5 * time.Second)
						continue
					}

					if appendErr := chain.Append(nextBlock); appendErr != nil {
						zap.L().Error("append block fail", zap.Error(appendErr))
						continue
					}

					nextCheckIndex++
				}
			}()

			<-ctx.Done()
		},
	}

	return cmd
}

func init() {
	dappCmd.AddCommand(NewDAppNode())
}
