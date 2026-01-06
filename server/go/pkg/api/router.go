package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/segolab/relay-ref/server/go/pkg/middleware"
	"github.com/segolab/relay-ref/server/go/pkg/ratelimit"
	"github.com/segolab/relay-ref/server/go/pkg/store"
)

type Config struct {
	HTTPAddr       string
	APIKeys        map[string]struct{}
	MaxBodyBytes   int64
	IdempotencyTTL time.Duration
	LimitPostRPS   float64
	LimitPostBurst int
	LimitGetRPS    float64
	LimitGetBurst  int
	LogLevel       slog.Level
}

type Dependencies struct {
	Logger      *slog.Logger
	Config      Config
	RelayStore  store.RelayStore
	Idempotency store.IdempotencyStore
	Limiter     ratelimit.Limiter
}

type App struct {
	Router http.Handler
}

func NewApp(d Dependencies) *App {
	h := NewHandlers(d.Logger, d.Config, d.RelayStore, d.Idempotency)

	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Recover(d.Logger))
	r.Use(middleware.RequestID())
	r.Use(middleware.RespondJSON())

	// System endpoints (no auth)
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	r.Get("/readyz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })

	// API group
	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.APIKeyAuth(d.Config.APIKeys))
		// Rate limiting by route-group (keeps diagrams clean and matches “per route” policy)
		r.With(middleware.RateLimit(d.Limiter, "post_relays")).
			Post("/relays", h.CreateRelay)
		r.With(middleware.RateLimit(d.Limiter, "get_relays")).
			Get("/relays", h.ListRelays)
		r.With(middleware.RateLimit(d.Limiter, "get_relays")).
			Get("/relays/{id}", h.GetRelay)
	})

	return &App{Router: r}
}
