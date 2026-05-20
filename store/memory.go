package store

import (
	"sync"

	"github.com/HtetAungKhant23/shortix/shortener"
)

type InMemoryStore struct {
	mu       sync.RWMutex
	entities map[string]*shortener.URL
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		entities: make(map[string]*shortener.URL),
	}
}

func (m *InMemoryStore) Save(url *shortener.URL) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.entities[url.ShortCode]; exists {
		return shortener.ErrCodeExists
	}

	cp := *url
	m.entities[url.ShortCode] = &cp
	return nil
}
