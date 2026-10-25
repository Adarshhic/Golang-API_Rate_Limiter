
package main

import (
	"fmt"
	"net"
	"net/http"
)

func RateLimitMiddleware(
	next http.Handler,
	limiter *RateLimiter,
) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Invalid client address", http.StatusBadRequest)
			return
		}

		if !limiter.Allow(ip) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintln(w, "Rate limit exceeded. Please try again later.")
			return
		}

		next.ServeHTTP(w, r)
	})
}
