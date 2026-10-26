package main

import (
	"errors"
	"sync"
	"time"
)

type TokenBucket struct {
	tokens     float64
	capacity   int
	rate       float64
	lastRefill time.Time
	lastSeen   time.Time
}

type AtomicBucketStore interface {
	Allow(ip string, rate float64, capacity int) (bool, error)
}

type RateLimiter struct {
	store    BucketStore
	atomic   AtomicBucketStore
	mu       sync.Mutex
	rate     float64
	capacity int
}

func NewRateLimiter(store BucketStore, rate float64, capacity int) (*RateLimiter, error) {
	if store == nil {
		return nil, errors.New("bucket store cannot be nil")
	}

	if rate <= 0 {
		return nil, errors.New("rate must be greater than zero")
	}

	if capacity <= 0 {
		return nil, errors.New("capacity must be greater than zero")
	}

	limiter := &RateLimiter{
		store:    store,
		rate:     rate,
		capacity: capacity,
	}

	if atomicStore, ok := store.(AtomicBucketStore); ok {
		limiter.atomic = atomicStore
	}

	return limiter, nil
}

func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()

	tb.tokens += elapsed * tb.rate

	if tb.tokens > float64(tb.capacity) {
		tb.tokens = float64(tb.capacity)
	}

	tb.lastRefill = now
}

func (rl *RateLimiter) Allow(ip string) bool {
	allowed, err := rl.AllowWithError(ip)
	return err == nil && allowed
}

func (rl *RateLimiter) AllowWithError(ip string) (bool, error) {
	// Redis supports atomic rate limiting.
	if rl.atomic != nil {
		return rl.atomic.Allow(ip, rl.rate, rl.capacity)
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	bucket, exists := rl.store.Get(ip)

	if !exists {
		bucket = &TokenBucket{
			tokens:     float64(rl.capacity),
			capacity:   rl.capacity,
			rate:       rl.rate,
			lastRefill: now,
			lastSeen:   now,
		}
	} else {
		bucket.refill()
		bucket.lastSeen = now
	}

	allowed := bucket.tokens >= 1

	if allowed {
		bucket.tokens--
	}

	rl.store.Set(ip, bucket)

	return allowed, nil
}

func (rl *RateLimiter) CleanupInactiveClients(timeout time.Duration) {
	rl.store.DeleteInactive(timeout)
}

func (rl *RateLimiter) StartCleanup(
	interval time.Duration,
	timeout time.Duration,
) {
	if interval <= 0 || timeout <= 0 {
		return
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			rl.CleanupInactiveClients(timeout)
		}
	}()
}

var ErrInvalidLimiterConfig = errors.New("rate and capacity must be positive")
