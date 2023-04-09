package cmd

import (
	"context"
	"fmt"
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
		transactionTopic    = "boatswain/transaction"
		reconciliationTopic = "boatswain/reconciliation"

		discoveryServiceTag        = "github.com/goforbroke1006/boatswain/dapp/chat"
		allInterfacesAnyFreePortMA = "/ip4/0.0.0.0/tcp/0"

		historyTailLength = 1024
	)

	var (
		nickArg = chat.DefaultNick()
		roomArg = "awesome-chat-room"
	)

	cmd := &cobra.Command{
		Use: "chat",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			// create a new p2p Host that listens on a random TCP port
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

			// TODO: validate room name
			chatTopic := fmt.Sprintf("chat/%s", roomArg)

			msgStreamIn, msgStreamInErr := messaging.NewStreamIn[*domain.TransactionPayload](
				ctx, chatTopic, p2pPubSub, p2pHost.ID(), false)
			if msgStreamInErr != nil {
				zap.L().Fatal("fail", zap.Error(msgStreamInErr))
			}
			msgStreamOut, msgStreamOutErr := messaging.NewStreamOut[domain.TransactionPayload](
				ctx, chatTopic, p2pPubSub)
			if msgStreamOutErr != nil {
				zap.L().Fatal("fail", zap.Error(msgStreamOutErr))
			}

			txStreamOut, txStreamOutErr := messaging.NewStreamOut[domain.TransactionPayload](
				ctx, transactionTopic, p2pPubSub)
			if txStreamOutErr != nil {
				zap.L().Fatal("fail", zap.Error(txStreamOutErr))
			}

			reconStreamIn, reconStreamInErr := messaging.NewStreamIn[*domain.ReconciliationResp](
				ctx, reconciliationTopic, p2pPubSub, p2pHost.ID(), true)
			if reconStreamInErr != nil {
				zap.L().Fatal("fail", zap.Error(reconStreamInErr))
			}

			historyMixer := chat.NewHistoryMixer(historyTailLength, msgStreamIn.In(), reconStreamIn.In())
			go func() {
				if runErr := historyMixer.Run(ctx); runErr != nil {
					panic(runErr)
				}
			}()

			// draw the UI
			ui := chat.NewChatUI(
				p2pPubSub,
				nickArg, chatTopic,
				historyMixer,
				msgStreamOut.Out(),
				txStreamOut.Out())
			go func() {
				if runErr := ui.Run(ctx); runErr != nil {
					zap.L().Fatal("running text UI fail", zap.Error(runErr))
				}
			}()

			<-ctx.Done()
		},
	}

	// parse some flags to set our nickname and the room to join
	cmd.PersistentFlags().StringVar(&nickArg, "nick", nickArg, "nickname to use in chat. will be generated if empty")
	cmd.PersistentFlags().StringVar(&roomArg, "room", roomArg, "name of chat room to join")

	return cmd
}

func init() {
	dappCmd.AddCommand(NewDAppChat())
}
