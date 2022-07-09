package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/libp2p/go-libp2p-core/peer"
)

// printErr is like fmt.Printf, but writes to stderr.
func printErr(m string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, m, args...)
}

// defaultNick generates a nickname based on the $USER environment variable and
// the last 8 chars of a peer ID.
func defaultNick(p peer.ID) string {
	currentUser, _ := user.Current()
	return fmt.Sprintf("%s-%s", currentUser.Username, shortID(p))
}

// shortID returns the last 8 chars of a base58-encoded peer id.
func shortID(p peer.ID) string {
	pretty := p.Pretty()
	return pretty[len(pretty)-8:]
}
