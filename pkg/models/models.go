package models

// ChatMessage gets converted to/from JSON and sent in the body of pubsub messages.
type State struct {
	Blocks []Block
}
type Block struct {
	Id           int
	Hash         string
	PreviousHash string
	Data         string
	Nonce        int
}
