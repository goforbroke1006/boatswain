package cmd

import (
	"context"
	"os"
	"os/signal"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/goforbroke1006/boatswain/internal/component/util/chat"
)

func NewUtilChat() *cobra.Command {
	var (
		nickArg = chat.DefaultNick()
		roomArg = "awesome-chat-room"
	)

	cmd := &cobra.Command{
		Use: "chat",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()

			// setup local mDNS discovery
			discovery := chat.NewDiscovery(p2pHost, discoveryServiceTag)
			if discoveryErr := discovery.Start(); discoveryErr != nil {
				zap.L().Fatal("initialize gossip sub fail", zap.Error(discoveryErr))
			}
			defer func() { _ = discovery.Close() }()

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
	utilCmd.AddCommand(NewUtilChat())
}
