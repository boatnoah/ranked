package sortedsets

import "github.com/go-redis/redis/v8"

// NewRedisClient builds a Redis client with the provided connection details.
func NewRedisClient(addr, pw string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       db,
	})
}
