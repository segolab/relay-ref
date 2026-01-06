package api

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

func LoadConfigFromEnv() Config {
	return Config{
		HTTPAddr:       getenv("RELAY_HTTP_ADDR", ":8429"),
		APIKeys:        parseAPIKeys(getenv("RELAY_API_KEYS", "dev-key")),
		MaxBodyBytes:   int64(getenvInt("RELAY_MAX_BODY_BYTES", 32768)),
		IdempotencyTTL: time.Duration(getenvInt("RELAY_IDEMPOTENCY_TTL_SECONDS", 3600)) * time.Second,
		LimitPostRPS:   getenvFloat("RELAY_LIMIT_POST_RPS", 10),
		LimitPostBurst: getenvInt("RELAY_LIMIT_POST_BURST", 20),
		LimitGetRPS:    getenvFloat("RELAY_LIMIT_GET_RPS", 50),
		LimitGetBurst:  getenvInt("RELAY_LIMIT_GET_BURST", 100),
		LogLevel:       slog.LevelInfo,
	}
}

func getenv(k, def string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	return def
}

func getenvInt(k string, def int) int {
	v := strings.TrimSpace(os.Getenv(k))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func getenvFloat(k string, def float64) float64 {
	v := strings.TrimSpace(os.Getenv(k))
	if v == "" {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return def
	}
	return f
}

func parseAPIKeys(csv string) map[string]struct{} {
	out := map[string]struct{}{}
	for _, part := range strings.Split(csv, ",") {
		k := strings.TrimSpace(part)
		if k != "" {
			out[k] = struct{}{}
		}
	}
	return out
}
