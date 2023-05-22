package mdns

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

func NewNotifee(h host.Host) mdns.Notifee {
	return &discoveryNotifee{h: h}
}

var _ mdns.Notifee = &discoveryNotifee{}

// discoveryNotifee gets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
	h host.Host
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
func (n *discoveryNotifee) HandlePeerFound(addInfo peer.AddrInfo) {
	fmt.Printf("discovered new peer %s\n", addInfo.String())

	err := n.h.Connect(context.Background(), addInfo)
	if err != nil {
		fmt.Printf("error connecting to peer %s: %s\n", addInfo.ID.String(), err)
	}
}
