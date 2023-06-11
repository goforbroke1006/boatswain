package cmd

import (
	"context"
	"os"
	"os/signal"
	"os/user"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/goforbroke1006/boatswain/internal/component/dapp/chat"
	"github.com/goforbroke1006/boatswain/internal/component/dapp/chat/node_client"
)

func NewDAppChat() *cobra.Command {
	var (
		nodeAPIAddrArg         = "http://localhost:58687"
		userNameArg            = "noname"
		dhtRendezvousPhraseArg = "github.com/goforbroke1006/boatswain/chat"
	)

	currUser, currUserErr := user.Current()
	if currUserErr != nil {
		zap.L().Error("get username failed", zap.Error(currUserErr))
	}
	userNameArg = currUser.Username

	cmd := &cobra.Command{
		Use:   "chat",
		Short: "Chat sample",
		Run: func(cmd *cobra.Command, args []string) {
			if len(userNameArg) == 0 {
				zap.L().Error("username is required")
				os.Exit(1)
			}

			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			nodeClient, nodeClientErr := node_client.NewClientWithResponses(nodeAPIAddrArg)
			if nodeClientErr != nil {
				zap.L().Fatal("http client initialization failed", zap.Error(nodeClientErr))
			}

			// draw the UI
			ui := chat.NewChatUI(nodeClient, userNameArg)
			go func() {
				if runErr := ui.Run(ctx); runErr != nil {
					zap.L().Fatal("running text UI fail", zap.Error(runErr))
				}
				stop()
			}()

			<-ctx.Done()
		},
	}

	cmd.PersistentFlags().StringVar(&nodeAPIAddrArg, "node-addr", nodeAPIAddrArg, "Node API address")
	cmd.PersistentFlags().StringVar(&userNameArg, "username", userNameArg, "Chat nick name")
	cmd.PersistentFlags().StringVar(&dhtRendezvousPhraseArg, "rendezvous", dhtRendezvousPhraseArg,
		"DHT rendezvous phrase should be same for all peers in network")

	return cmd
}

func init() {
	dappCmd.AddCommand(NewDAppChat())
}
