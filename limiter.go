
package main

import (
		"sync"
		"time"
) 

type TokenBucket struct {
	tokens     float64
	capacity   int
	rate       float64
	lastRefill time.Time
}

type RateLimiter struct {
	clients map[string]*TokenBucket
	rate     float64
    capacity int
	mu      sync.Mutex
} 




func NewRateLimiter(rate float64, capacity int , ) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*TokenBucket),
		rate:    rate,
		capacity: capacity,
	}
}



func(tb *TokenBucket)refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	newTokens := elapsed * tb.rate
	tb.tokens +=newTokens
	if tb.tokens > float64(tb.capacity){
		tb.tokens = float64(tb.capacity)
	}
		tb.lastRefill = now
}


func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.clients[ip]

	if !exists {
		bucket = &TokenBucket{
			tokens:     float64(rl.capacity),
			capacity:   rl.capacity,
			rate:       rl.rate,
			lastRefill: time.Now(),
		}

		rl.clients[ip] = bucket
	}

	bucket.refill()

	if bucket.tokens >= 1 {
		bucket.tokens--
		return true
	}

	return false
}
