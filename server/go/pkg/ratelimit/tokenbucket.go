package ratelimit

import (
	"math"
	"sync"
	"time"
)

type Result struct {
	Allowed           bool
	Limit             int
	Remaining         int
	ResetInSeconds    int
	RetryAfterSeconds int
}

type Limiter interface {
	Allow(apiKey, routeGroup string, now time.Time) Result
}

type Config struct {
	PostRPS   float64
	PostBurst int
	GetRPS    float64
	GetBurst  int
}

type tokenBucket struct {
	rps   float64
	burst float64

	mu     sync.Mutex
	last   time.Time
	tokens float64
}

func newTokenBucket(rps float64, burst int) *tokenBucket {
	return &tokenBucket{
		rps:    rps,
		burst:  float64(burst),
		last:   time.Now().UTC(),
		tokens: float64(burst),
	}
}

func (b *tokenBucket) allow(now time.Time) Result {
	b.mu.Lock()
	defer b.mu.Unlock()

	dt := now.Sub(b.last).Seconds()
	if dt < 0 {
		dt = 0
	}
	b.tokens = math.Min(b.burst, b.tokens+(b.rps*dt))
	b.last = now

	limit := int(b.burst)
	remaining := int(math.Floor(b.tokens))
	if remaining < 0 {
		remaining = 0
	}

	if b.tokens >= 1.0 {
		b.tokens -= 1.0
		remaining = int(math.Floor(b.tokens))
		if remaining < 0 {
			remaining = 0
		}
		return Result{
			Allowed:           true,
			Limit:             limit,
			Remaining:         remaining,
			ResetInSeconds:    0,
			RetryAfterSeconds: 0,
		}
	}

	// Denied: estimate time until 1 token is available.
	need := 1.0 - b.tokens
	secs := int(math.Ceil(need / b.rps))
	if secs < 1 {
		secs = 1
	}
	return Result{
		Allowed:           false,
		Limit:             limit,
		Remaining:         remaining,
		ResetInSeconds:    secs,
		RetryAfterSeconds: secs,
	}
}

type TokenBucketLimiter struct {
	cfg Config

	mu   sync.Mutex
	post map[string]*tokenBucket
	get  map[string]*tokenBucket
}

func NewTokenBucketLimiter(cfg Config) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		cfg:  cfg,
		post: make(map[string]*tokenBucket),
		get:  make(map[string]*tokenBucket),
	}
}

func (l *TokenBucketLimiter) Allow(apiKey, routeGroup string, now time.Time) Result {
	if apiKey == "" {
		// Should not happen (auth runs before), but be safe.
		return Result{Allowed: false, Limit: 0, Remaining: 0, ResetInSeconds: 1, RetryAfterSeconds: 1}
	}

	switch routeGroup {
	case "post_relays":
		return l.allowFromMap(l.post, apiKey, l.cfg.PostRPS, l.cfg.PostBurst, now)
	default:
		return l.allowFromMap(l.get, apiKey, l.cfg.GetRPS, l.cfg.GetBurst, now)
	}
}

func (l *TokenBucketLimiter) allowFromMap(m map[string]*tokenBucket, apiKey string, rps float64, burst int, now time.Time) Result {
	l.mu.Lock()
	b := m[apiKey]
	if b == nil {
		b = newTokenBucket(rps, burst)
		m[apiKey] = b
	}
	l.mu.Unlock()
	return b.allow(now)
}
