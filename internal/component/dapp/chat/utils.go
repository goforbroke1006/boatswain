package chat

import (
	"fmt"
	"os/user"

	"github.com/libp2p/go-libp2p/core/peer"
)

// DefaultNick generates a nickname based on the $USER environment variable.
func DefaultNick() string {
	currentUser, _ := user.Current()
	return currentUser.Username
}

// shortID returns the last 8 chars of a base58-encoded peer id.
func shortID(p peer.ID) string {
	pretty := p.Pretty()
	return pretty[len(pretty)-8:]
}

// withColor wraps a string with color tags for display in the messages text box.
func withColor(color, msg string) string {
	return fmt.Sprintf("[%s]%s[-]", color, msg)
}
