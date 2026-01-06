package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segolab/relay-ref/server/go/pkg/api"
	"github.com/segolab/relay-ref/server/go/pkg/ratelimit"
	"github.com/segolab/relay-ref/server/go/pkg/store"
)

func main() {
	cfg := api.LoadConfigFromEnv()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))

	relayStore := store.NewInMemoryRelayStore()
	idem := store.NewInMemoryIdempotencyStore(cfg.IdempotencyTTL)

	limiter := ratelimit.NewTokenBucketLimiter(ratelimit.Config{
		PostRPS:   cfg.LimitPostRPS,
		PostBurst: cfg.LimitPostBurst,
		GetRPS:    cfg.LimitGetRPS,
		GetBurst:  cfg.LimitGetBurst,
	})

	app := api.NewApp(api.Dependencies{
		Logger:      logger,
		Config:      cfg,
		RelayStore:  relayStore,
		Idempotency: idem,
		Limiter:     limiter,
	})

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           app.Router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("server started", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen failed", "err", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("shutting down")
	_ = srv.Close()
}
