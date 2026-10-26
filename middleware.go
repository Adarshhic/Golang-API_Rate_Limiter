
package main

import (
	"net"
	"net/http"
)


func RateLimitMiddleware(
	next http.Handler,
	limiter *RateLimiter,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the IP address from RemoteAddr (IP:port).
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(
				w,
				"Invalid client address",
				http.StatusBadRequest,
			)
			return
		}

	
		allowed, err := limiter.AllowWithError(ip)
		if err != nil {
			http.Error(
				w,
				"Rate limiter temporarily unavailable",
				http.StatusServiceUnavailable,
			)
			return
		}

	
		if !allowed {
			http.Error(
				w,
				"Rate limit exceeded. Please try again later.",
				http.StatusTooManyRequests,
			)
			return
		}


		next.ServeHTTP(w, r)
	})
}
