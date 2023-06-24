package cmd

import (
	"context"
	"os/signal"
	"syscall"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/goforbroke1006/boatswain/internal/common"
	"github.com/goforbroke1006/boatswain/internal/system"
	"github.com/goforbroke1006/boatswain/pkg/discovery/discovery_dht"
)

func NewDAppChat() *cobra.Command {
	var (
		handleMultiAddrArg     = "/ip4/0.0.0.0/tcp/58688"
		userNameArg            = system.MustGetCurrentUsername()
		dhtRendezvousPhraseArg = "github.com/goforbroke1006/boatswain/dapp/chat/v1"
	)

	cmd := &cobra.Command{
		Use:     "chat",
		Version: "v1.0",
		Short:   "Chat sample",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			// load key pair or create
			privateKey, publicKey, keysPairErr := common.GetKeysPair()
			if keysPairErr != nil {
				zap.L().Fatal("get keys pair failed", zap.Error(keysPairErr))
			}

			// start peer
			p2pHost, p2pHostErr := libp2p.New(
				libp2p.Identity(privateKey),
				libp2p.ListenAddrStrings(handleMultiAddrArg),
			)
			if p2pHostErr != nil {
				zap.L().Fatal("p2p host listening fail", zap.Error(p2pHostErr))
			}
			defer func() { _ = p2pHost.Close() }()
			zap.L().Info("host peer started", zap.String("id", p2pHost.ID().String()))

			// connect with another peers with bootstraps
			bootstrapPeers, bootstrapPeersErr := common.LoadPeersList()
			if bootstrapPeersErr != nil {
				zap.L().Fatal("load bootstrap peers list failed", zap.Error(bootstrapPeersErr))
			}
			dht, dhtErr := discovery_dht.NewKDHT(ctx, p2pHost, bootstrapPeers)
			if dhtErr != nil {
				zap.L().Fatal("initialize DHT failed", zap.Error(dhtErr))
			}
			zap.L().Info("bootstrap peers connecting", zap.Int("count", len(bootstrapPeers)))
			go discovery_dht.Discover(ctx, p2pHost, dht, dhtRendezvousPhraseArg)
			zap.L().Info("discovery initialized")

			// create a new PubSub service using the GossipSub router
			p2pPubSub, p2pPubSubErr := pubsub.NewGossipSub(ctx, p2pHost)
			if p2pPubSubErr != nil {
				zap.L().Fatal("initialize gossip sub fail", zap.Error(p2pPubSubErr))
			}
			zap.L().Info("pub-sub initialized")

			_ = privateKey
			_ = publicKey
			_ = p2pPubSub

			// draw the UI

			a := app.New()
			w := a.NewWindow("Hello")

			hello := widget.NewLabel("Hello Fyne!")
			w.SetContent(container.NewVBox(
				hello,
				widget.NewButton("Hi!", func() {
					hello.SetText("Welcome :)")
				}),
			))

			w.ShowAndRun()

			//app := chat.NewApplication(userNameArg, privateKey, publicKey, p2pPubSub)
			//_ = app
			//go func() {
			//	if runErr := app.Run(ctx); runErr != nil && runErr != context.Canceled {
			//		zap.L().Fatal("app running fail", zap.Error(runErr))
			//	}
			//	zap.L().Debug("exit UI")
			//	stop()
			//}()

			//<-ctx.Done()
		},
	}

	cmd.PersistentFlags().StringVar(&handleMultiAddrArg, "addr", handleMultiAddrArg,
		"Host listen this multi-address")
	cmd.PersistentFlags().StringVar(&userNameArg, "username", userNameArg, "Chat nick name")
	cmd.PersistentFlags().StringVar(&dhtRendezvousPhraseArg, "rendezvous", dhtRendezvousPhraseArg,
		"DHT rendezvous phrase should be same for all peers in network")

	return cmd
}

func init() {
	dappCmd.AddCommand(NewDAppChat())
}
