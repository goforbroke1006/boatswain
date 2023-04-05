package cmd

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"time"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/goforbroke1006/boatswain/internal/blockchain"
	"github.com/goforbroke1006/boatswain/internal/component/chat"
	"github.com/goforbroke1006/boatswain/internal/storage"
)

func NewDAppChat() *cobra.Command {
	const (
		// discoveryServiceTag is used in our mDNS advertisements to discover other chat peers.
		discoveryServiceTag        = "boatswain-chat-example"
		allInterfacesAnyFreePortMA = "/ip4/0.0.0.0/tcp/0"
	)

	var (
		nickArg = chat.DefaultNick()
		roomArg = "awesome-chat-room"
	)

	cmd := &cobra.Command{
		Use: "chat",
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

			db, err := sql.Open("sqlite3", "./chat-blocks.db")
			if err != nil {
				panic(err)
			}
			schemaQuery, err := os.ReadFile("./schema.sql")
			if err != nil {
				panic(err)
			}
			if _, err := db.Exec(string(schemaQuery)); err != nil {
				panic(err)
			}
			blockStorage := storage.NewBlockStorage(db)
			chain := blockchain.NewBlockChain(blockStorage)

			messageCb := func(message pubsub.Message) {
				_ = chain.Generate(time.Now().Unix(), string(message.Data))
			}
			// join the chat room
			cr, err := chat.JoinChatRoom(ctx, p2pPubSub, p2pHost.ID(), nickArg, roomArg, messageCb)
			if err != nil {
				zap.L().Fatal("join chat room fail", zap.Error(err))
			}

			// draw the UI
			ui := chat.NewChatUI(cr)
			if err = ui.Run(); err != nil {
				zap.L().Fatal("running text UI fail", zap.Error(err))
			}
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
