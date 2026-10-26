
package main

import "testing"

  
func TestRateLimiter(t *testing.T) {
	store := NewMemoryStore()
	limiter, err := NewRateLimiter(store, 0.001, 2)
if err != nil {
    t.Fatal(err)
}

	// Bucket starts with 2 tokens.
	if !limiter.Allow("192.168.1.1") {
		t.Fatal("first request should be allowed")
	}

	if !limiter.Allow("192.168.1.1") {
		t.Fatal("second request should be allowed")
	}

	// Bucket is empty; no time to refill.
	if limiter.Allow("192.168.1.1") {
		t.Fatal("third request should be rejected")
	}
}

func TestDifferentClients(t *testing.T) {
	store := NewMemoryStore()
	limiter, err := NewRateLimiter(store, 2, 5)
if err != nil {
	t.Fatal(err)
}

	if !limiter.Allow("192.168.1.1") {
		t.Fatal("first client should be allowed")
	}

	// Each IP has its own bucket.
	if !limiter.Allow("192.168.1.2") {
		t.Fatal("second client should have its own bucket")
	}
}
