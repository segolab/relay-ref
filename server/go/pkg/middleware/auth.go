package middleware

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey string

const apiKeyCtxKey ctxKey = "api_key"

func APIKeyAuth(allowed map[string]struct{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			k := strings.TrimSpace(r.Header.Get("X-API-Key"))
			if k == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if _, ok := allowed[k]; !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), apiKeyCtxKey, k)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func APIKeyFromContext(ctx context.Context) string {
	v := ctx.Value(apiKeyCtxKey)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
