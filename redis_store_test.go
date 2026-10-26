
package main

import (
	"fmt"
	"testing"
	"time"
)

func TestRedisStoreAllow(t *testing.T) {
	store := NewRedisStore("localhost:6379")

	if err := store.client.Ping(store.ctx).Err(); err != nil {
		t.Skipf("Redis is unavailable: %v", err)
	}

	ip := fmt.Sprintf("test-client-%d", time.Now().UnixNano())
	key := "rate_limit:" + ip

	defer store.client.Del(store.ctx, key)

	// Capacity = 2, refill rate = 1 token per second.
	allowed, err := store.Allow(ip, 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if !allowed {
		t.Fatal("first request should be allowed")
	}

	allowed, err = store.Allow(ip, 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if !allowed {
		t.Fatal("second request should be allowed")
	}

	allowed, err = store.Allow(ip, 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if allowed {
		t.Fatal("third request should be rejected")
	}
}
