package cmd

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/goforbroke1006/boatswain/domain"
	"github.com/goforbroke1006/boatswain/internal/component/dapp/chat"
	"github.com/goforbroke1006/boatswain/pkg/discovery"
	"github.com/goforbroke1006/boatswain/pkg/messaging"
)

func NewDAppChat() *cobra.Command {
	const (
		chatTopic               = "boatswain/dapp/chat"
		transactionTopic        = "boatswain/_transaction"
		reconciliationRespTopic = "boatswain/_reconciliation/resp"

		discoveryServiceTag        = "github.com/goforbroke1006/boatswain/dapp/chat"
		allInterfacesAnyFreePortMA = "/ip4/0.0.0.0/tcp/0"

		historyTailLength = 1024
	)

	cmd := &cobra.Command{
		Use: "chat",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			p2pHost, p2pHostErr := libp2p.New(libp2p.ListenAddrStrings(allInterfacesAnyFreePortMA))
			if p2pHostErr != nil {
				zap.L().Fatal("p2p host listening fail", zap.Error(p2pHostErr))
			}
			defer func() { _ = p2pHost.Close() }()
			zap.L().Info("host peer started",
				zap.String("peer-id", p2pHost.ID().String()),
				zap.Any("addresses", p2pHost.Addrs()))

			nickName := chat.DefaultNick(p2pHost)

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

			msgStream, msgStreamErr := messaging.NewStreamBoth[
				domain.Transaction,
				*domain.Transaction,
			](ctx, chatTopic, p2pPubSub, p2pHost.ID(), false)
			if msgStreamErr != nil {
				zap.L().Fatal("fail", zap.Error(msgStreamErr))
			}

			txStreamOut, txStreamOutErr := messaging.NewStreamOut[domain.Transaction](
				ctx, transactionTopic, p2pPubSub)
			if txStreamOutErr != nil {
				zap.L().Fatal("fail", zap.Error(txStreamOutErr))
			}

			reconStreamIn, reconStreamInErr := messaging.NewStreamIn[
				domain.ReconciliationResp,
				*domain.ReconciliationResp,
			](ctx, reconciliationRespTopic, p2pPubSub, p2pHost.ID(), true)
			if reconStreamInErr != nil {
				zap.L().Fatal("fail", zap.Error(reconStreamInErr))
			}

			historyMixer := chat.NewHistoryMixer(historyTailLength, msgStream.In(), reconStreamIn.In())
			go func() {
				if runErr := historyMixer.Run(ctx); runErr != nil {
					panic(runErr)
				}
			}()

			// draw the UI
			ui := chat.NewChatUI(
				p2pHost, p2pPubSub, chatTopic,
				nickName,
				historyMixer,
				msgStream.Out(),
				txStreamOut.Out())
			go func() {
				if runErr := ui.Run(ctx); runErr != nil {
					zap.L().Fatal("running text UI fail", zap.Error(runErr))
				}
			}()

			<-ctx.Done()
		},
	}

	return cmd
}

func init() {
	dappCmd.AddCommand(NewDAppChat())
}
