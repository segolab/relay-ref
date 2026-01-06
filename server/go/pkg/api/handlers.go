package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/segolab/relay-ref/server/go/pkg/middleware"
	"github.com/segolab/relay-ref/server/go/pkg/model"
	"github.com/segolab/relay-ref/server/go/pkg/store"
)

type Handlers struct {
	log   *slog.Logger
	cfg   Config
	store store.RelayStore
	idem  store.IdempotencyStore
}

func NewHandlers(log *slog.Logger, cfg Config, s store.RelayStore, idem store.IdempotencyStore) *Handlers {
	return &Handlers{log: log, cfg: cfg, store: s, idem: idem}
}

func (h *Handlers) CreateRelay(w http.ResponseWriter, r *http.Request) {
	apiKey := middleware.APIKeyFromContext(r.Context())
	if apiKey == "" {
		WriteError(w, r, http.StatusUnauthorized, "unauthorized", "missing API key", nil)
		return
	}

	// Bounded read
	r.Body = http.MaxBytesReader(w, r.Body, h.cfg.MaxBodyBytes)
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "invalid_request", "failed to read request body", map[string]any{"err": err.Error()})
		return
	}

	var req model.CreateRelayRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		WriteError(w, r, http.StatusBadRequest, "invalid_request", "invalid JSON", map[string]any{"err": err.Error()})
		return
	}
	if err := validateCreate(req); err != nil {
		WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error(), nil)
		return
	}
	if !isJSONObject(req.Payload) {
		WriteError(w, r, http.StatusBadRequest, "invalid_request", "payload must be a JSON object", nil)
		return
	}

	payloadHash := sha256.Sum256(raw)
	hashHex := hex.EncodeToString(payloadHash[:])

	idemKey := strings.TrimSpace(r.Header.Get("Idempotency-Key"))

	createFn := func() (*model.Relay, error) {
		now := time.Now().UTC()
		relay := &model.Relay{
			ID:            uuid.New(),
			EventType:     req.EventType,
			Destination:   req.Destination,
			Payload:       req.Payload,
			Metadata:      req.Metadata,
			Status:        model.RelayStatusQueued,
			CreatedAt:     now,
			DeliveredAt:   nil,
			FailureReason: nil,
		}
		h.store.Create(relay)
		return relay, nil
	}

	var relay *model.Relay

	if idemKey != "" {
		res, err := h.idem.GetOrCreate(apiKey, idemKey, hashHex, createFn)
		if err != nil {
			if store.IsIdempotencyConflict(err) {
				WriteError(w, r, http.StatusConflict, "idempotency_conflict", "idempotency key reuse with different payload", nil)
				return
			}
			WriteError(w, r, http.StatusInternalServerError, "internal", "idempotency failure", map[string]any{"err": err.Error()})
			return
		}
		relay = res
	} else {
		relay, err = createFn()
		if err != nil {
			WriteError(w, r, http.StatusInternalServerError, "internal", "failed to create relay", map[string]any{"err": err.Error()})
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(relay)
}

func (h *Handlers) GetRelay(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "invalid_request", "invalid relay id", nil)
		return
	}

	relay, ok := h.store.Get(id)
	if !ok {
		WriteError(w, r, http.StatusNotFound, "not_found", "relay not found", nil)
		return
	}

	_ = json.NewEncoder(w).Encode(relay)
}

func (h *Handlers) ListRelays(w http.ResponseWriter, r *http.Request) {
	pageSize := 50
	if v := r.URL.Query().Get("pageSize"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 1 && n <= 200 {
			pageSize = n
		}
	}

	pageToken := r.URL.Query().Get("pageToken")
	offset, err := store.DecodePageToken(pageToken)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "invalid_request", "invalid pageToken", nil)
		return
	}

	items, nextOffset := h.store.List(pageSize, offset)

	resp := model.ListRelaysResponse{
		Items:         items,
		NextPageToken: store.EncodePageToken(nextOffset),
	}
	// If no more items, return null in JSON (matches spec's nullable)
	if nextOffset < 0 {
		resp.NextPageToken = nil
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func validateCreate(req model.CreateRelayRequest) error {
	if strings.TrimSpace(req.EventType) == "" || len(req.EventType) > 128 {
		return errf("eventType is required and must be <= 128 characters")
	}
	if req.Destination.Type != "webhook" {
		return errf("destination.type must be 'webhook'")
	}
	if strings.TrimSpace(req.Destination.URL) == "" || len(req.Destination.URL) > 2048 {
		return errf("destination.url is required and must be <= 2048 characters")
	}
	if len(req.Payload) == 0 {
		return errf("payload is required")
	}
	return nil
}

func isJSONObject(raw json.RawMessage) bool {
	// Minimal check: first non-space must be '{'
	for _, b := range raw {
		if b == ' ' || b == '\n' || b == '\r' || b == '\t' {
			continue
		}
		return b == '{'
	}
	return false
}

type simpleErr string

func (e simpleErr) Error() string { return string(e) }
func errf(msg string) error       { return simpleErr(msg) }
