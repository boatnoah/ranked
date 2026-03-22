package sortedsets

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

const allTimeKey = "lb:alltime"

// Entry represents a leaderboard row in Redis.
type Entry struct {
	UserID   int64 `json:"user_id"`
	Trophies int64 `json:"trophies"`
	Rank     int64 `json:"rank"`
}

type Store interface {
	Incr(ctx context.Context, userID int64, delta int64) (int64, error)
	Top(ctx context.Context, limit int64) ([]Entry, error)
	Rank(ctx context.Context, userID int64) (Entry, error)
	BulkSet(ctx context.Context, entries []Entry) error
	Reset(ctx context.Context) error
}

type RedisStore struct {
	cmd redis.Cmdable
}

func NewRedisStore(cmd redis.Cmdable) *RedisStore {
	return &RedisStore{cmd: cmd}
}

func (s *RedisStore) Incr(ctx context.Context, userID int64, delta int64) (int64, error) {
	member := strconv.FormatInt(userID, 10)
	oldScore, err := s.cmd.ZScore(ctx, allTimeKey, member).Result()

	if err != nil {
		return 0, err
	}
	newScore := max(0, int64(oldScore)+delta)

	err = s.cmd.ZAdd(ctx, allTimeKey, &redis.Z{Member: member, Score: float64(newScore)}).Err()

	if err != nil {
		return 0, err
	}

	return newScore, nil
}

func (s *RedisStore) Top(ctx context.Context, limit int64) ([]Entry, error) {
	if limit <= 0 {
		return []Entry{}, nil
	}

	zs, err := s.cmd.ZRevRangeWithScores(ctx, allTimeKey, 0, limit-1).Result()
	if err != nil {
		return nil, err
	}

	out := make([]Entry, 0, len(zs))
	for i, z := range zs {
		memberStr, ok := z.Member.(string)
		if !ok {
			return nil, redis.Nil
		}
		userID, err := strconv.ParseInt(memberStr, 10, 64)
		if err != nil {
			return nil, err
		}

		out = append(out, Entry{
			UserID:   userID,
			Trophies: int64(z.Score),
			Rank:     int64(i + 1),
		})
	}

	return out, nil
}

func (s *RedisStore) Rank(ctx context.Context, userID int64) (Entry, error) {
	member := strconv.FormatInt(userID, 10)

	rank, err := s.cmd.ZRevRank(ctx, allTimeKey, member).Result()
	if err != nil {
		return Entry{}, err
	}

	score, err := s.cmd.ZScore(ctx, allTimeKey, member).Result()
	if err != nil {
		return Entry{}, err
	}

	return Entry{
		UserID:   userID,
		Trophies: int64(score),
		Rank:     rank + 1, // convert to 1-based
	}, nil
}

func (s *RedisStore) BulkSet(ctx context.Context, entries []Entry) error {
	if len(entries) == 0 {
		return nil
	}

	zs := make([]*redis.Z, 0, len(entries))
	for _, e := range entries {
		zs = append(zs, &redis.Z{
			Score:  float64(e.Trophies),
			Member: strconv.FormatInt(e.UserID, 10),
		})
	}

	return s.cmd.ZAdd(ctx, allTimeKey, zs...).Err()
}

func (s *RedisStore) Reset(ctx context.Context) error {
	return s.cmd.Del(ctx, allTimeKey).Err()
}
