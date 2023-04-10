package chat

import (
	"fmt"
	"os/user"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

// DefaultNick generates a nickName based on the $USER environment variable.
func DefaultNick(p2pHost host.Host) string {
	currentUser, _ := user.Current()
	return fmt.Sprintf("%s - %s", currentUser.Username, p2pHost.ID().String())
}

// shortID returns the last 8 chars of a base58-encoded peer id.
func shortID(p peer.ID) string {
	pretty := p.String()
	return pretty[len(pretty)-8:]
}

// withColor wraps a string with color tags for display in the messages text box.
func withColor(color, msg string) string {
	return fmt.Sprintf("[%s]%s[-]", color, msg)
}
