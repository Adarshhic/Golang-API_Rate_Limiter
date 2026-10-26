
package main

import (
	"sync"
	"time"
)


type MemoryStore struct {
	clients map[string]*TokenBucket
	mu      sync.Mutex
}


func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		clients: make(map[string]*TokenBucket),
	}
}


func (s *MemoryStore) Get(ip string) (*TokenBucket, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	bucket, exists := s.clients[ip]
	return bucket, exists
}


func (s *MemoryStore) Set(ip string, bucket *TokenBucket) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[ip] = bucket
}


func (s *MemoryStore) DeleteInactive(timeout time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := time.Now().Add(-timeout)

	for ip, bucket := range s.clients {
		if bucket.lastSeen.Before(cutoff) {
			delete(s.clients, ip)
		}
	}
}
