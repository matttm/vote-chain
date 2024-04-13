package messenger

import (
	"context"
	"encoding/json"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"

	"vote-chain/pkg/models"
	"vote-chain/pkg/topics"
	"vote-chain/internal/utilities"
)

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
	sub        *pubsub.Subscription

	roomName string
	self     peer.ID
	nick     string
}

// type stateMessage struct {
// 	Payload models.State
// 	SenderID string
// }
func CreateMessenger() *Messenger {
	m := new(Messenger)
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
func (msngr *Messenger) ReadLoop() {
	for {
		msg, err := msngr.sub.Next(msngr.ctx)
		if err != nil {
			close(msngr.StateChannel)
			panic("Error occured reading sub")
		}
		// only forward messages delivered by others
		if msg.ReceivedFrom == msngr.self {
			panic("Error resding message from self")
		}
		stateMessage := new(models.State)
		err = json.Unmarshal(msg.Data, stateMessage)
		if err != nil {
			continue
		}
		// send valid messages onto the Messages channel
		msngr.StateChannel <- stateMessage
	}
}


func (msngr *Messenger) ListenToVoteChain(ctx context.Context, ps *pubsub.PubSub, selfID peer.ID) {
	// join the pubsub topic
	topic := utilities.CreateTopicHandle(ps, topics.CHAIN)

	// and subscribe to it
	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}
	// start reading messages from the subscription in a loop
	// ToDo:  put this kon its own thread
	go msngr.ReadLoop()
}
