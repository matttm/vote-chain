package models

// ChatMessage gets converted to/from JSON and sent in the body of pubsub messages.
type State struct {
	blocks []Block
}
type Block struct {
	Id int
	Hash string
	PreviousHash ztring
	Datan string
	Nonce int
}
