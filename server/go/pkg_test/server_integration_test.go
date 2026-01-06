package pkg_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"log/slog"

	"github.com/segolab/relay-ref/server/go/pkg/api"
	"github.com/segolab/relay-ref/server/go/pkg/ratelimit"
	"github.com/segolab/relay-ref/server/go/pkg/store"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	cfg := api.Config{
		HTTPAddr:       ":0",
		APIKeys:        map[string]struct{}{"k": {}},
		MaxBodyBytes:   32768,
		IdempotencyTTL: 1 * time.Hour,
		LimitPostRPS:   1,
		LimitPostBurst: 2, // burst > 1 is required to test idempotency: retries must still respect rate limits but not be blocked immediately
		LimitGetRPS:    50,
		LimitGetBurst:  100,
		LogLevel:       slog.LevelInfo,
	}

	app := api.NewApp(api.Dependencies{
		Logger:      slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)),
		Config:      cfg,
		RelayStore:  store.NewInMemoryRelayStore(),
		Idempotency: store.NewInMemoryIdempotencyStore(cfg.IdempotencyTTL),
		Limiter: ratelimit.NewTokenBucketLimiter(ratelimit.Config{
			PostRPS:   cfg.LimitPostRPS,
			PostBurst: cfg.LimitPostBurst,
			GetRPS:    50,
			GetBurst:  100,
		}),
	})
	return httptest.NewServer(app.Router)
}

func TestUnauthorized(t *testing.T) {
	s := newTestServer(t)
	defer s.Close()

	reqBody := []byte(`{"eventType":"x","destination":{"type":"webhook","url":"https://e"},"payload":{"a":1}}`)
	resp, err := http.Post(s.URL+"/v1/relays", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestCreateAndGet(t *testing.T) {
	s := newTestServer(t)
	defer s.Close()

	body := map[string]any{
		"eventType": "order.created",
		"destination": map[string]any{
			"type": "webhook",
			"url":  "https://example.com/hook",
		},
		"payload": map[string]any{"x": 1},
	}
	raw, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", s.URL+"/v1/relays", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "k")
	req.Header.Set("Idempotency-Key", "idem1")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	var created map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&created)
	id := created["id"].(string)

	getReq, _ := http.NewRequest("GET", s.URL+"/v1/relays/"+id, nil)
	getReq.Header.Set("X-API-Key", "k")
	getResp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatal(err)
	}
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", getResp.StatusCode)
	}
}

func TestIdempotencySameKeyReturnsSameRelay(t *testing.T) {
	s := newTestServer(t)
	defer s.Close()

	raw := []byte(`{"eventType":"x","destination":{"type":"webhook","url":"https://e"},"payload":{"a":1}}`)

	do := func() string {
		req, _ := http.NewRequest("POST", s.URL+"/v1/relays", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", "k")
		req.Header.Set("Idempotency-Key", "idem42")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected 201, got %d", resp.StatusCode)
		}
		var m map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&m)
		return m["id"].(string)
	}

	id1 := do()
	id2 := do()
	if id1 != id2 {
		t.Fatalf("expected same relay id, got %s and %s", id1, id2)
	}
}

func TestRateLimit429(t *testing.T) {
	s := newTestServer(t)
	defer s.Close()

	raw := []byte(`{"eventType":"x","destination":{"type":"webhook","url":"https://e"},"payload":{"a":1}}`)

	req := func() int {
		r, _ := http.NewRequest("POST", s.URL+"/v1/relays", bytes.NewReader(raw))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("X-API-Key", "k")
		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			t.Fatal(err)
		}
		return resp.StatusCode
	}

	_ = req() // consumes burst token
	code := req()
	if code != http.StatusTooManyRequests {
		// Depending on timing, second may pass. Third should fail.
		code = req()
	}
	if code != http.StatusTooManyRequests {
		t.Fatalf("expected at least one 429, got %d", code)
	}
}
