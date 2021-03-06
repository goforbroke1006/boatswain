package main

import (
	"context"
	"database/sql"
	"flag"
	"io/ioutil"
	"time"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	_ "github.com/mattn/go-sqlite3"

	"github.com/goforbroke1006/boatswain/internal/blockchain"
	"github.com/goforbroke1006/boatswain/internal/storage"
)

// discoveryServiceTag is used in our mDNS advertisements to discover other chat peers.
const discoveryServiceTag = "pubsub-chat-example"

const allInterfacesAnyFreePortMA = "/ip4/0.0.0.0/tcp/0"

func main() {
	// parse some flags to set our nickname and the room to join
	nickFlag := flag.String("nick", "", "nickname to use in chat. will be generated if empty")
	roomFlag := flag.String("room", "awesome-chat-room", "name of chat room to join")
	flag.Parse()

	ctx := context.Background()

	// create a new libp2p Host that listens on a random TCP port
	h, err := libp2p.New(libp2p.ListenAddrStrings(allInterfacesAnyFreePortMA))
	if err != nil {
		panic(err)
	}

	// create a new PubSub service using the GossipSub router
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		panic(err)
	}

	// setup local mDNS discovery
	discovery := NewDiscovery(h, discoveryServiceTag)
	if err := discovery.Start(); err != nil {
		panic(err)
	}
	defer func() { _ = discovery.Close() }()

	// use the nickname from the cli flag, or a default if blank
	nick := *nickFlag
	if len(nick) == 0 {
		nick = defaultNick(h.ID())
	}

	// join the room from the cli flag, or the flag default
	room := *roomFlag

	db, err := sql.Open("sqlite3", "./chat-blocks.db")
	if err != nil {
		panic(err)
	}
	schemaQuery, err := ioutil.ReadFile("./schema.sql")
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(string(schemaQuery)); err != nil {
		panic(err)
	}
	blockStorage := storage.NewBlockStorage(db)
	chain := blockchain.NewBlockChain(blockStorage)

	messageCb := func(message pubsub.Message) {
		_ = chain.Generate(time.Now().Unix(), string(message.Data))
	}
	// join the chat room
	cr, err := JoinChatRoom(ctx, ps, h.ID(), nick, room, messageCb)
	if err != nil {
		panic(err)
	}

	// draw the UI
	ui := NewChatUI(cr)
	if err = ui.Run(); err != nil {
		printErr("error running text UI: %s", err)
	}
}
