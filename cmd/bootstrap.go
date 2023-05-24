package cmd

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/goforbroke1006/boatswain/internal/common"
	"github.com/goforbroke1006/boatswain/pkg/discovery/discovery_dht"
	"github.com/goforbroke1006/go-healthcheck"
)

func NewBootstrap() *cobra.Command {
	const (
		allInterfacesCertainPortMA = "/ip4/0.0.0.0/tcp/9999"
	)

	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstrap service",
		Long:  "Bootstrap service required to help nodes discover each other",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			healthcheck.Panel().Start(ctx, healthcheck.DefaultAddr)

			privateKey, privateKeyErr := common.ReadPrivateKey()
			if privateKeyErr != nil {
				zap.L().Fatal("read private key failed", zap.Error(privateKeyErr))
			}

			p2pHost, p2pHostErr := libp2p.New(
				libp2p.Identity(privateKey),
				libp2p.ListenAddrStrings(allInterfacesCertainPortMA),
			)
			if p2pHostErr != nil {
				zap.L().Fatal("p2p host listening fail", zap.Error(p2pHostErr))
			}
			defer func() { _ = p2pHost.Close() }()
			zap.L().Info("host peer started",
				zap.String("peer-id", p2pHost.ID().String()),
				zap.Any("addresses", p2pHost.Addrs()))

			log.Printf("Connect to me on:")
			for _, addr := range p2pHost.Addrs() {
				log.Printf("  %s/p2p/%s", addr, p2pHost.ID().String())
			}

			dht, dhtErr := discovery_dht.NewKDHT(ctx, p2pHost, nil)
			if dhtErr != nil {
				zap.L().Fatal("initialize DHT discovery fail", zap.Error(dhtErr))
			}

			healthcheck.Panel().SetHealthy()

			go discovery_dht.Discover(ctx, p2pHost, dht, DHTRendezvousPhrase)

			healthcheck.Panel().SetReady()

			<-ctx.Done()
		},
	}

	return cmd
}

func init() {
	rootCmd.AddCommand(NewBootstrap())
}
