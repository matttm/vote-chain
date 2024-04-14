package controllers

import (
	"vote-chain/internal/delegate"

	"github.com/go-chi/chi/v5"
)

type BallotController struct {
	delegate *delegate.Delegate
}

func CreateBallotController(delegate *delegate.Delegate) *BallotController {
	b := new(BallotController)
	b.delegate = delegate
	return b
}
func (b *BallotController) Router() chi.Router {
	r := chi.NewRouter()
	return r
}
