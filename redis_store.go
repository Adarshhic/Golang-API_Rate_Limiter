
package main

import (
	"context"
	"time"
   "fmt"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStore(addr string) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisStore{
		client: client,
		ctx:    context.Background(),
	}
}

// Get implements BucketStore.
// Redis buckets are stored as hashes by the atomic script.
func (s *RedisStore) Get(ip string) (*TokenBucket, bool) {
	key := "rate_limit:" + ip

	values, err := s.client.HMGet(
		s.ctx,
		key,
		"tokens",
		"last_refill",
	).Result()

	if err != nil || values[0] == nil || values[1] == nil {
		return nil, false
	}

	tokens, err1 := parseFloat(values[0])
	lastRefill, err2 := parseFloat(values[1])

	if err1 != nil || err2 != nil {
		return nil, false
	}

	return &TokenBucket{
		tokens:     tokens,
		capacity:   0,
		rate:       0,
		lastRefill: time.Unix(0, int64(lastRefill*1e9)),
	}, true
}

// Set implements BucketStore.
func (s *RedisStore) Set(ip string, bucket *TokenBucket) {
	// Atomic Redis.Allow manages the hash directly.
	// This method is retained to satisfy BucketStore.
}

// DeleteInactive implements BucketStore.
// Redis expires each bucket automatically in the Lua script.
func (s *RedisStore) DeleteInactive(timeout time.Duration) {
	// No action required: bucket keys have a Redis TTL.
}

// Allow performs the token-bucket operation atomically in Redis.
func (s *RedisStore) Allow(
	ip string,
	rate float64,
	capacity int,
) (bool, error) {
	script := redis.NewScript(`
		local key = KEYS[1]
		local capacity = tonumber(ARGV[1])
		local rate = tonumber(ARGV[2])

		local time_parts = redis.call("TIME")
		local now = tonumber(time_parts[1])
			+ tonumber(time_parts[2]) / 1000000

		local data = redis.call("HMGET", key, "tokens", "last_refill")
		local tokens = tonumber(data[1])
		local last_refill = tonumber(data[2])

		if tokens == nil then
			tokens = capacity
			last_refill = now
		end

		local elapsed = math.max(0, now - last_refill)
		tokens = math.min(capacity, tokens + elapsed * rate)

		local allowed = 0
		if tokens >= 1 then
			tokens = tokens - 1
			allowed = 1
		end

		redis.call("HSET", key,
			"tokens", tokens,
			"last_refill", now
		)
		redis.call("EXPIRE", key, 600)

		return allowed
	`)

	result, err := script.Run(
		s.ctx,
		s.client,
		[]string{"rate_limit:" + ip},
		capacity,
		rate,
	).Int()

	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func parseFloat(value interface{}) (float64, error) {
	switch v := value.(type) {
	case string:
		var result float64
		_, err := fmt.Sscan(v, &result)
		return result, err
	case []byte:
		var result float64
		_, err := fmt.Sscan(string(v), &result)
		return result, err
	default:
		return 0, fmt.Errorf("unexpected Redis value type %T", value)
	}
}
