package utilities

import (
	"log"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func CreateTopicHandle(ps *pubsub.PubSub, topic string) *pubsub.Topic {
	handle, err := ps.Join(topic)
	if err != nil {
		log.Panic("Error: ", err)
	}
	return handle
}
