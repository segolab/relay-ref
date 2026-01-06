package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"

const requestIDCtxKey ctxKey = "request_id"

func RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(requestIDHeader)
			if id == "" {
				id = uuid.NewString()
			}
			w.Header().Set(requestIDHeader, id)
			ctx := context.WithValue(r.Context(), requestIDCtxKey, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequestIDFromContext(ctx context.Context) *string {
	v := ctx.Value(requestIDCtxKey)
	if s, ok := v.(string); ok && s != "" {
		return &s
	}
	return nil
}
