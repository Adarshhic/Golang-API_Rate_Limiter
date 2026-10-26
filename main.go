package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func main() {

	store := NewRedisStore("localhost:6379")

	if err := store.client.Ping(store.ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	limiter, err := NewRateLimiter(store, 2, 5)
	if err != nil {
		log.Fatalf("Invalid rate limiter configuration: %v", err)
	}

	handler := RateLimitMiddleware(
		http.HandlerFunc(helloHandler),
		limiter,
	)

	http.Handle("/", handler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
