package cmd

import (
	"context"
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
		nodeAPIAddrArg = "http://localhost:8081"
		userNameArg    = "noname"
	)

	currUser, currUserErr := user.Current()
	if currUserErr != nil {
		zap.L().Error("get username failed", zap.Error(currUserErr))
	}
	userNameArg = currUser.Username

	cmd := &cobra.Command{
		Use: "chat",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			client, clientErr := node_client.NewClientWithResponses(nodeAPIAddrArg)
			if clientErr != nil {
				zap.L().Fatal("http client initialization failed", zap.Error(clientErr))
			}

			// draw the UI
			ui := chat.NewChatUI(client, userNameArg)
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

	return cmd
}

func init() {
	dappCmd.AddCommand(NewDAppChat())
}
