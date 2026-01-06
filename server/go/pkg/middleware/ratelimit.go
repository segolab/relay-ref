package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/segolab/relay-ref/server/go/pkg/ratelimit"
)

func RateLimit(l ratelimit.Limiter, routeGroup string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := APIKeyFromContext(r.Context())
			res := l.Allow(apiKey, routeGroup, time.Now().UTC())

			// Headers on best-effort basis
			w.Header().Set("RateLimit-Limit", strconv.Itoa(res.Limit))
			w.Header().Set("RateLimit-Remaining", strconv.Itoa(res.Remaining))
			w.Header().Set("RateLimit-Reset", strconv.Itoa(res.ResetInSeconds))

			if !res.Allowed {
				w.Header().Set("Retry-After", strconv.Itoa(res.RetryAfterSeconds))
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
