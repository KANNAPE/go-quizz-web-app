package lobby

import (
	"sync"

	"github.com/google/uuid"
)

type memory struct {
	mu   sync.RWMutex
	data map[string]Lobby
}

func NewMemoryStore() Store {
	return &memory{data: map[string]Lobby{}}
}

func (m *memory) Create(host string) (Lobby, error) {
	id := uuid.New().String()
	l := Lobby{ID: id, HostName: host}
	m.mu.Lock()
	m.data[id] = l
	m.mu.Unlock()
	return l, nil
}

func (m *memory) Get(id string) (Lobby, bool) {
	m.mu.RLock()
	l, ok := m.data[id]
	m.mu.RUnlock()
	return l, ok
}
