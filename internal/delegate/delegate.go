package delegate

import (
	"vote-chain/internal/messenger"
	"vote-chain/pkg/models"
)

type Delegate struct {
	messenger *messenger.Messenger
	state     *models.State
}

func CreateDelegate() *Delegate {
	d := new(Delegate)
	d.messenger = messenger.CreateMessenger()
	d.state = nil

	d.messenger.ListenToVoteChain()
	return d
}
func (d *Delegate) DestroyDelegate() {}
