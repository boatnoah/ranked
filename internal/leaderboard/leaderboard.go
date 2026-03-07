package leaderboard

import (
	"errors"

	"github.com/boatnoah/ranked/internal/storage"
)

type Leaderboard struct {
	storage storage.Storage
}

type MatchPayload struct {
	Result string
	Crowns int32
}

func (l *Leaderboard) Submit(matchPayload MatchPayload) error {

	err := validatePayload(matchPayload)
	if err != nil {

		return err
	}

	return nil

}

func validatePayload(matchPayload MatchPayload) error {
	if matchPayload.Crowns < 0 || matchPayload.Crowns > 3 {
		return errors.New("Valid crowns must be in the range of 0 - 3")
	}

	if matchPayload.Result != "win" && matchPayload.Result != "loss" {
		return errors.New("Result must be a win or a loss")
	}
	return nil

}
