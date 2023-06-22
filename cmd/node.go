package cmd

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/goforbroke1006/go-healthcheck"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/goforbroke1006/boatswain/domain"
	"github.com/goforbroke1006/boatswain/internal"
	"github.com/goforbroke1006/boatswain/internal/common"
	"github.com/goforbroke1006/boatswain/internal/component/node"
	"github.com/goforbroke1006/boatswain/internal/component/node/api/impl"
	"github.com/goforbroke1006/boatswain/internal/component/node/api/spec"
	"github.com/goforbroke1006/boatswain/pkg/discovery/discovery_dht"
)

func NewNode() *cobra.Command {
	var (
		handleMultiAddrArg     = "/ip4/0.0.0.0/tcp/58687"
		dhtRendezvousPhraseArg = "github.com/goforbroke1006/boatswain"
	)

	cmd := &cobra.Command{
		Use:   "node",
		Short: "Node component",
		Long:  "Node appends transactions and sync blocks",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			healthcheck.Panel().Start(ctx, healthcheck.DefaultAddr)
			healthcheck.Panel().SetHealthy()

			// load key pair or create
			privateKey, privateKeyErr := common.ReadPrivateKey()
			if privateKeyErr != nil {
				zap.L().Fatal("read private key failed", zap.Error(privateKeyErr))
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

			var (
				txCache         = internal.NewTransactionCache()
				txSpreadSvc     = internal.NewTransactionSpreadInfoService(p2pPubSub)
				txReader        domain.TransactionReader
				blockSpreadSvc  domain.BlockVoteSpreadInfoService
				blockStorage    domain.BlockStorage
				blockVoteReader domain.BlockVoteReader
				voteCollector   domain.VoteCollector
			)

			const requiredForBlockCount = 12 // TODO: move to config

			app := node.NewApplication(
				p2pHost,
				txCache,
				txReader,
				blockStorage,
				blockSpreadSvc,
				requiredForBlockCount,
				blockVoteReader,
				voteCollector)
			go func() {
				if runErr := app.Run(ctx); runErr != nil && runErr != context.Canceled {
					zap.L().Fatal("run application failed", zap.Error(runErr))
				}
			}()

			router := echo.New()
			router.HideBanner = true
			router.Use(middleware.Recover())
			router.Use(middleware.CORS())
			spec.RegisterHandlers(router, impl.NewHandlers(p2pHost, txCache, txSpreadSvc))
			go func() {
				if startErr := router.Start("0.0.0.0:58687"); startErr != nil {
					zap.L().Fatal("start http server failed", zap.Error(startErr))
				}
			}()

			healthcheck.Panel().SetReady()

			<-ctx.Done()
		},
	}

	cmd.PersistentFlags().StringVar(&handleMultiAddrArg, "addr", handleMultiAddrArg,
		"Host listen this multi-address")
	cmd.PersistentFlags().StringVar(&dhtRendezvousPhraseArg, "rendezvous", dhtRendezvousPhraseArg,
		"DHT rendezvous phrase should be same for all peers in network")

	return cmd
}

func init() {
	rootCmd.AddCommand(NewNode())
}
