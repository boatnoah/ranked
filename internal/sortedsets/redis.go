package sortedsets

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(addr, pw string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       db,
	})
}

// functions
// Add Score
// Retrieve Top Player

func AddPlayer(ctx context.Context, client *redis.Client, score float64, email string) error {
	if email == "" {
		return errors.New("Invalid email")
	}

	member := redis.Z{
		Score:  score,
		Member: email,
	}

	client.ZAdd(ctx, "leaderboard", &member)

	return nil
}

func GetUserRank(ctx context.Context, client *redis.Client, email string) (int64, error) {
	rank, err := client.ZRevRank(ctx, "leaderboard", email).Result()

	if err != nil {
		return 0, err
	}

	return rank, nil
}

func GetTopPlayers(ctx context.Context, client *redis.Client) ([]redis.Z, error) {
	leaderboard, err := client.ZRevRangeWithScores(ctx, "leaderboard", 0, 2).Result()

	if err != nil {
		return nil, err
	}

	return leaderboard, nil

}
