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

func (m *InMemoryStore) FindByCode(code string) (*shortener.URL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entity, ok := m.entities[code]
	if !ok {
		return nil, shortener.ErrNotFound
	}

	cp := *entity

	return &cp, nil
}

func (m *InMemoryStore) IncrementAccess(code string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	entity, ok := m.entities[code]
	if !ok {
		return shortener.ErrNotFound
	}

	entity.AccessCount++
	return nil
}

func (m *InMemoryStore) Update(code string, newURL string) (*shortener.URL, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	entity, ok := m.entities[code]
	if !ok {
		return nil, shortener.ErrNotFound
	}

	entity.URL = newURL
	entity.AccessCount = 0

	cp := *entity
	return &cp, nil
}

func (m *InMemoryStore) Delete(code string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.entities[code]; !ok {
		return shortener.ErrNotFound
	}

	delete(m.entities, code)

	return nil
}
