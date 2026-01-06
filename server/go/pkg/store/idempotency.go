package store

import (
	"errors"
	"sync"
	"time"

	"github.com/segolab/relay-ref/server/go/pkg/model"
)

type IdempotencyStore interface {
	GetOrCreate(apiKey, idemKey, payloadHash string, createFn func() (*model.Relay, error)) (*model.Relay, error)
}

type idemEntry struct {
	payloadHash string
	relay       *model.Relay
	expiresAt   time.Time
}

type InMemoryIdempotencyStore struct {
	ttl time.Duration

	mu sync.Mutex
	m  map[string]idemEntry
}

func NewInMemoryIdempotencyStore(ttl time.Duration) *InMemoryIdempotencyStore {
	return &InMemoryIdempotencyStore{
		ttl: ttl,
		m:   make(map[string]idemEntry),
	}
}

var errIdemConflict = errors.New("idempotency conflict")

func IsIdempotencyConflict(err error) bool {
	return errors.Is(err, errIdemConflict)
}

func (s *InMemoryIdempotencyStore) GetOrCreate(apiKey, idemKey, payloadHash string, createFn func() (*model.Relay, error)) (*model.Relay, error) {
	now := time.Now().UTC()
	k := apiKey + ":" + idemKey

	s.mu.Lock()
	// Cleanup opportunistically
	if e, ok := s.m[k]; ok {
		if now.After(e.expiresAt) {
			delete(s.m, k)
		} else {
			// Found
			if e.payloadHash != payloadHash {
				s.mu.Unlock()
				return nil, errIdemConflict
			}
			relay := e.relay
			s.mu.Unlock()
			return relay, nil
		}
	}
	s.mu.Unlock()

	relay, err := createFn()
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.m[k] = idemEntry{
		payloadHash: payloadHash,
		relay:       relay,
		expiresAt:   now.Add(s.ttl),
	}
	s.mu.Unlock()

	return relay, nil
}
