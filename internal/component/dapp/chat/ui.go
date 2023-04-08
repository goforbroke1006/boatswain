package chat

import (
	"context"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/rivo/tview"

	"github.com/goforbroke1006/boatswain/domain"
)

// NewChatUI returns a new ChatUI struct that controls the text UI.
// It won't actually do anything until you call Run().
func NewChatUI(
	p2pPubSub *pubsub.PubSub,
	nickname string,
	chatTopic string,
	historyMixer *HistoryMixer,
	msgOut chan<- *domain.TransactionPayload,
	txOut chan<- *domain.TransactionPayload,
) *ChatUI {
	return &ChatUI{
		p2pPubSub:    p2pPubSub,
		nickname:     nickname,
		chatTopic:    chatTopic,
		historyMixer: historyMixer,
		msgOut:       msgOut,
		txOut:        txOut,
	}
}

// ChatUI is a Text User Interface (TUI) for a ChatRoom.
// The Run method will draw the UI to the terminal in "fullscreen"
// mode. You can quit with Ctrl-C, or by typing "/quit" into the
// chat prompt.
type ChatUI struct {
	p2pPubSub *pubsub.PubSub

	nickname  string
	chatTopic string

	historyMixer *HistoryMixer

	msgOut chan<- *domain.TransactionPayload
	txOut  chan<- *domain.TransactionPayload
}

// Run starts the chat event loop in the background, then starts
// the event loop for the text UI.
func (ui *ChatUI) Run(ctx context.Context) error {
	app := tview.NewApplication()

	// make a text view to contain our chat messages
	msgBox := tview.NewTextView()
	msgBox.SetDynamicColors(true)
	msgBox.SetBorder(true)
	msgBox.SetTitle(fmt.Sprintf("Room: %s", ui.chatTopic))
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
				history, _ := ui.historyMixer.History()
				msgBox.Clear()
				for _, item := range history {
					color := "red"
					if item.PeerSender == ui.nickname {
						color = "green"
					}
					prompt := withColor(color, fmt.Sprintf("<%s>:", item.PeerSender))
					_, _ = fmt.Fprintf(msgBox, "%s %s\n", prompt, item.Content)
				}
			}
		}
	}()

	// text views are io.Writers, but they don't automatically refresh.
	// this sets a change handler to force the app to redraw when we get
	// new messages to display.
	msgBox.SetChangedFunc(func() { app.Draw() })

	// an input field for typing messages into
	input := tview.NewInputField().
		SetLabel(ui.nickname + " > ").
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
		prompt := withColor("green", fmt.Sprintf("<%s>:", ui.nickname))
		_, _ = fmt.Fprintf(msgBox, "%s %s\n", prompt, line)

		// TODO: send message to room mates
		// TODO: send message to node

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
			case <-time.After(10 * time.Second):
				peers := ui.p2pPubSub.ListPeers(ui.chatTopic)
				peersList.Clear()
				for _, p := range peers {
					_, _ = fmt.Fprintln(peersList, shortID(p))
				}

				app.Draw()
			}
		}
	}()

	// chatPanel is a horizontal box with messages on the left and peers on the right
	// the peers list takes 20 columns, and the messages take the remaining space
	chatPanel := tview.NewFlex().
		AddItem(msgBox, 0, 1, false).
		AddItem(peersList, 20, 1, false)

	// flex is a vertical box with the chatPanel on top and the input field at the bottom.

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatPanel, 0, 1, false).
		AddItem(input, 1, 1, true)

	app.SetRoot(flex, true)

	errsCh := make(chan error, 2)

	go func() {
		if appRunErr := app.Run(); appRunErr != nil {
			errsCh <- appRunErr
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errsCh:
		return err
	}
}
