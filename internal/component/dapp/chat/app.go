package chat

import (
	"context"
	"crypto"
	"fmt"
	"go.uber.org/zap"
	"time"

	"github.com/gdamore/tcell/v2"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/rivo/tview"

	"github.com/goforbroke1006/boatswain/internal/component/dapp/chat/node_client"
)

// NewApplication returns a new Application struct that controls the text UI.
// It won't actually do anything until you call Init().
func NewApplication(
	nickName string,
	privateKey crypto.PrivateKey,
	publicKey crypto.PublicKey,
	p2pPubSub *pubsub.PubSub,
) *Application {
	return &Application{nickName: nickName}
}

// Application is a Text User Interface (TUI) for a ChatRoom.
// The Run method will draw the UI to the terminal in "fullscreen"
// mode. You can quit with Ctrl-C, or by typing "/quit" into the
// chat prompt.
type Application struct {
	client   node_client.ClientWithResponsesInterface
	nickName string
}

// Run starts the chat event loop in the background, then starts
// the event loop for the text UI.
func (ui *Application) Run(ctx context.Context) error {
	zap.L().Debug("starting UI")

	app := tview.NewApplication()
	defer app.Stop()

	// make a text view to contain our chat messages
	msgBox := tview.NewTextView()
	msgBox.SetDynamicColors(true)
	msgBox.SetBorder(true)
	msgBox.SetTitle("Chat Room")

	// text views are io.Writers, but they don't automatically refresh.
	// this sets a change handler to force the app to redraw when we get
	// new messages to display.
	msgBox.SetChangedFunc(func() { app.Draw() })

	// an input field for typing messages into
	input := tview.NewInputField().
		SetLabel(ui.nickName + " > ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorBlack)
	// the done func is called when the user hits enter, or tabs out of the field
	input.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			// we don't want to do anything if they just tabbed away
			return
		}
		line := input.GetText()
		if len(line) == 0 {
			// ignore blank lines
			return
		}

		// bail if requested
		if line == "/quit" {
			app.Stop()
			return
		}

		// display message in UI
		prompt := withColor("green", fmt.Sprintf("<%s>:", ui.nickName))
		_, _ = fmt.Fprintf(msgBox, "%s %s\n", prompt, line)

		// TODO: send message to roommates

		input.SetText("")
	})

	// make a text view to hold the list of peers in the room, updated by ui.refreshPeers()
	peersList := tview.NewTextView()
	peersList.SetBorder(true)
	peersList.SetTitle("Peers")
	peersList.SetChangedFunc(func() { app.Draw() })
	go func() { // refresh peers list
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
				// TODO: show address book

				app.Draw()
			}
		}
	}()

	// chatPanel is a horizontal box with messages on the left and peers on the right
	// the peers list takes 20 columns, and the messages take the remaining space
	chatPanel := tview.NewFlex().
		AddItem(msgBox, 0, 1, false).
		AddItem(peersList, 54, 1, false)

	// flex is a vertical box with the chatPanel on top and the input field at the bottom.

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatPanel, 0, 1, false).
		AddItem(input, 1, 1, true)

	app.SetRoot(flex, true)

	errsCh := make(chan error, 2)

	zap.L().Debug("running UI")
	go func() {
		if appRunErr := app.Run(); appRunErr != nil {
			errsCh <- appRunErr
		}
	}()

	zap.L().Debug("waiting for exit UI")
	select {
	case <-ctx.Done():
		zap.L().Debug("stopping UI")
		return ctx.Err()
	case err := <-errsCh:
		return err
	}
}
