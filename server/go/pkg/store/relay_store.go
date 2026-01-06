package store

import (
	"sync"

	"github.com/google/uuid"
	"github.com/segolab/relay-ref/server/go/pkg/model"
)

type RelayStore interface {
	Create(r *model.Relay)
	Get(id uuid.UUID) (*model.Relay, bool)
	List(pageSize int, offset int) (items []*model.Relay, nextOffset int)
}

type InMemoryRelayStore struct {
	mu    sync.RWMutex
	byID  map[uuid.UUID]*model.Relay
	order []uuid.UUID
}

func NewInMemoryRelayStore() *InMemoryRelayStore {
	return &InMemoryRelayStore{
		byID: make(map[uuid.UUID]*model.Relay),
	}
}

func (s *InMemoryRelayStore) Create(r *model.Relay) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byID[r.ID] = r
	s.order = append(s.order, r.ID)
}

func (s *InMemoryRelayStore) Get(id uuid.UUID) (*model.Relay, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.byID[id]
	return r, ok
}

func (s *InMemoryRelayStore) List(pageSize int, offset int) ([]*model.Relay, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if offset < 0 {
		offset = 0
	}
	if offset >= len(s.order) {
		return []*model.Relay{}, -1
	}

	end := offset + pageSize
	if end > len(s.order) {
		end = len(s.order)
	}

	out := make([]*model.Relay, 0, end-offset)
	for _, id := range s.order[offset:end] {
		out = append(out, s.byID[id])
	}

	if end >= len(s.order) {
		return out, -1
	}
	return out, end
}
