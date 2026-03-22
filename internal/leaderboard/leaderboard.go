package leaderboard

import (
	"context"
	"errors"
	"math/rand/v2"

	"github.com/boatnoah/ranked/internal/sortedsets"
	"github.com/boatnoah/ranked/internal/storage"
)

type Leaderboard struct {
	storage      *storage.Storage
	redisStorage *sortedsets.RedisStore
}

type MatchPayload struct {
	UserID int64
	Result string
	Crowns int64
}

func New(store *storage.Storage, redisStore *sortedsets.RedisStore) *Leaderboard {
	return &Leaderboard{store, redisStore}
}

func (l *Leaderboard) Submit(ctx context.Context, mp MatchPayload) (int64, error) {

	err := validatePayload(mp)
	if err != nil {
		return 0, err
	}

	delta := calcDelta(mp)

	var score int64

	err = l.storage.WithTx(ctx, func(ts storage.TxStorage) error {
		err = ts.MatchStore.Create(ctx, mp.UserID, mp.Result, mp.Crowns, delta)

		if err != nil {
			return err
		}

		err = ts.TrophyStore.Upsert(ctx, mp.UserID, delta)

		if err != nil {
			return err
		}

		s, err := l.redisStorage.Incr(ctx, mp.UserID, delta)
		if err != nil {
			return err
		}
		score = s
		return nil
	})

	if err != nil {
		return 0, err
	}

	return score, nil
}

func (l *Leaderboard) GetTopPlayers(ctx context.Context, limit int64) (*[]sortedsets.Entry, error) {
	entries, err := l.redisStorage.Top(ctx, limit)

	if err != nil {
		return nil, err
	}

	return &entries, nil

}

func (l *Leaderboard) GetPlayerRank(ctx context.Context, userID int64) (*sortedsets.Entry, error) {
	entry, err := l.redisStorage.Rank(ctx, userID)
	if err != nil {
		users, err := l.storage.TrophyStore.GetAllTrophies(ctx)

		if len(users) == 0 {
			return nil, errors.New("No players in the leaderboard")
		}

		var entries []sortedsets.Entry

		for _, user := range users {
			entry := sortedsets.Entry{UserID: user.UserID, Trophies: user.Trophies}

			entries = append(entries, entry)
		}

		l.redisStorage.BulkSet(ctx, entries)

		entry, err = l.redisStorage.Rank(ctx, userID)

		if err != nil {
			return nil, err
		}
	}

	return &entry, nil
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
	if matchPayload.Crowns == 3 && matchPayload.Result == "loss" {
		return errors.New("Breaking game logic; cannot lose with max crowns")
	}
	return nil

}

// bias = how many crowns / 3
// (1-bias)*r + bias (26 --- 34)
// (1-bias)*r

func calcDelta(matchPayload MatchPayload) int64 {
	crowns := matchPayload.Crowns

	r := rand.Float64()
	bias := float64(crowns) / 3.0

	switch matchPayload.Result {

	// adjust offset a but a reasonably amount

	case "win":
		mn := 26
		x := (1-bias)*r + bias
		res := int64(mn + (8 * int(x)))
		return res
	case "loss":
		mn := 22
		x := (1 - bias) * r
		res := int64(mn + (10 * int(x)))
		return -res
	default:
		return 0
	}
}
