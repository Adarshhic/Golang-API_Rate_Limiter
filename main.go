
package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received:", r.URL.Path)

		next.ServeHTTP(w, r)
	})
}


func main() {
	limiter := NewRateLimiter(2, 5)

	handler := RateLimitMiddleware(
		http.HandlerFunc(helloHandler),
		limiter,
	)

	http.Handle("/", handler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
