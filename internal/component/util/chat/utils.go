package chat

import (
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
