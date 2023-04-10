package chat

import (
	"fmt"
	"os/user"

	"github.com/libp2p/go-libp2p/core/host"
)

// DefaultNick generates a nickName based on the $USER environment variable.
func DefaultNick(p2pHost host.Host) string {
	currentUser, _ := user.Current()
	return fmt.Sprintf("%s - %s", currentUser.Username, p2pHost.ID().String())
}

// withColor wraps a string with color tags for display in the messages text box.
func withColor(color, msg string) string {
	return fmt.Sprintf("[%s]%s[-]", color, msg)
}
