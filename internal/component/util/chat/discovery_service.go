package chat

import (
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

// NewDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func NewDiscovery(h host.Host, tag string) mdns.Service {
	notifee := NewNotifee(h)
	s := mdns.NewMdnsService(h, tag, notifee) // setup mDNS discovery to find local peers
	return s
}
