package leaderboard

import (
	"context"
	"errors"
	"math/rand/v2"

	"github.com/boatnoah/ranked/internal/sortedsets"
	"github.com/boatnoah/ranked/internal/storage"
)

type Leaderboard struct {
	storage      storage.Storage
	redisStorage sortedsets.RedisStore
}

type MatchPayload struct {
	UserID int64
	Result string
	Crowns int64
}

const MAXCROWNS = 3

func (l *Leaderboard) Submit(ctx context.Context, matchPayload MatchPayload) error {

	err := validatePayload(matchPayload)
	if err != nil {
		return err
	}

	delta := calcDelta(matchPayload)

	l.storage.MatchStore.Create(ctx, matchPayload.UserID, matchPayload.Result, matchPayload.Crowns, delta)
	return nil

}

func validatePayload(matchPayload MatchPayload) error {
	if matchPayload.Crowns < 0 || matchPayload.Crowns > 3 {
		return errors.New("Valid crowns must be in the range of 0 - 3")
	}

	if matchPayload.Result != "win" && matchPayload.Result != "loss" {
		return errors.New("Result must be a win or a loss")
	}
	if matchPayload.Crowns == 0 && matchPayload.Result == "win" {
		return errors.New("Breaking game logic; cannot win with zero crowns")
	}
	if matchPayload.Crowns == MAXCROWNS && matchPayload.Result == "loss" {
		return errors.New("Breaking game logic; cannot loss with max crowns")
	}
	return nil

}

// bias = how many crowns / 3
// (1-bias)*r + bias (26 --- 34)
// (1-bias)*r

func calcDelta(matchPayload MatchPayload) int64 {
	crowns := matchPayload.Crowns

	r := rand.Int64()
	bias := crowns / MAXCROWNS

	switch matchPayload.Result {

	case "win":
		mn := 26
		x := (1-bias)*r + bias
		res := int64(mn + (10 * int(x)))
		return res
	case "loss":
		mn := 22
		x := (1 - bias) * r
		res := int64(mn + (8 * int(x)))
		return res
	default:
		return 0
	}
}
