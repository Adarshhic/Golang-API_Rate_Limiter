package main

import "testing"

func TestNewRateLimiterInvalidConfig(t *testing.T) {
	store := NewMemoryStore()

	tests := []struct {
		name     string
		rate     float64
		capacity int
	}{
		{"zero rate", 0, 5},
		{"negative rate", -1, 5},
		{"zero capacity", 2, 0},
		{"negative capacity", 2, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiter, err := NewRateLimiter(
				store,
				tt.rate,
				tt.capacity,
			)

			if err == nil {
				t.Fatal("expected an error for invalid configuration")
			}

			if limiter != nil {
				t.Fatal("expected limiter to be nil")
			}
		})
	}
}
