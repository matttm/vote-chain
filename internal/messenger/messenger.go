package messenger

import (
	"fmt"
	"time"

	"context"
	"encoding/json"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"

	"vote-chain/internal/utilities"
	"vote-chain/pkg/models"
	"vote-chain/pkg/topics"
)

// DiscoveryInterval is how often we re-publish our mDNS records.
const DiscoveryInterval = time.Hour

// ChatRoomBufSize is the number of incoming messages to buffer for each topic.
const ChatRoomBufSize = 128

// ChatRoom represents a subscription to a single PubSub topic. Messages
// can be published to the topic with ChatRoom.Publish, and received
// messages are pushed to the Messages channel.
type Messenger struct {
	// Messages is a channel of messages received from other peers in the chat room
	StateChannel chan *models.State

	ctx        context.Context
	ps         *pubsub.PubSub
	chainTopic *pubsub.Topic
	blockTopic *pubsub.Topic
	chainSub   *pubsub.Subscription
	self       peer.ID
}

// discoveryNotifee gets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
	h host.Host
}

//	type stateMessage struct {
//		Payload models.State
//		SenderID string
//	}
func CreateMessenger() *Messenger {
	m := new(Messenger)
	m.ctx = context.Background()

	protocolId := "/ip4/0.0.0.0/tcp/0"
	// create a new libp2p Host that listens on a random TCP port
	h, err := libp2p.New(libp2p.ListenAddrStrings(protocolId))
	if err != nil {
		panic(err)
	}
	floodsub, err := pubsub.NewFloodSub(m.ctx, h)
	if err != nil {
		panic(err)
	}
	m.ps = floodsub

	// setup local mDNS discovery
	if err := setupDiscovery(h); err != nil {
		panic(err)
	}
	return m
}

// Publish sends a message to the pubsub topic.
func (msngr *Messenger) Publish(message string) error {
	// TODO fill block
	m := models.Block{}
	msgBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if msngr.blockTopic == nil {
		msngr.blockTopic = utilities.CreateTopicHandle(msngr.ps, topics.CHAIN)
	}
	return msngr.chainTopic.Publish(msngr.ctx, msgBytes)
}

func (msngr *Messenger) ListPeers() []peer.ID {
	return msngr.ps.ListPeers(topics.CHAIN)
}

// readLoop pulls messages from the pubsub topic and pushes them onto the Messages channel.
func (msngr *Messenger) readLoop() {
	for {
		msg, err := msngr.chainSub.Next(msngr.ctx)
		if err != nil {
			close(msngr.StateChannel)
			panic("Error occured reading sub")
		}
		// only forward messages delivered by others
		if msg.ReceivedFrom == msngr.self {
			panic("Error resding message from self")
		}
		println("Received message")
		stateMessage := new(models.State)
		err = json.Unmarshal(msg.Data, stateMessage)
		if err != nil {
			continue
		}
		// send valid messages onto the Messages channel
		msngr.StateChannel <- stateMessage
	}
}

func (msngr *Messenger) ListenToVoteChain() {
	// join the pubsub topic
	topic := utilities.CreateTopicHandle(msngr.ps, topics.CHAIN)
	// and subscribe to it
	sub, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}
	msngr.chainSub = sub
	// start reading messages from the subscription in a loop
	// ToDo:  put this kon its own thread
	println("Listening to vote chain")
	go msngr.readLoop()
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Printf("discovered new peer %s\n", pi.ID)
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		fmt.Printf("error connecting to peer %s: %s\n", pi.ID, err)
	}
}

// setupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func setupDiscovery(h host.Host) error {
	// setup mDNS discovery to find local peers
	s := mdns.NewMdnsService(h, "", &discoveryNotifee{h: h})
	return s.Start()
}
