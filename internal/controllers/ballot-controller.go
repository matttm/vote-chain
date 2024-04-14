package controllers

import (
	"encoding/json"
	"net/http"
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
	r.Get("/", b.getBallotHandler)
	return r
}

func (b *BallotController) getBallotHandler(w http.ResponseWriter, req *http.Request) {
	println("Executing getBsllot request hsndler")
	data := b.delegate.State
	finalResponseJson, _ := json.Marshal(data)
	w.Write(finalResponseJson)
}
