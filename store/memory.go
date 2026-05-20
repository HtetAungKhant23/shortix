package store

import (
	"sync"

	"github.com/HtetAungKhant23/shortix/shortener"
)

type InMemoryStore struct {
	mu       *sync.RWMutex
	entities map[string]*shortener.URL
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		entities: make(map[string]*shortener.URL),
	}
}
