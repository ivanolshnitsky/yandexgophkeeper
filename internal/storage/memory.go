package storage

import "sync"

// MemoryStorage in-memory хранилище
type MemoryStorage struct {
	mu    sync.RWMutex
	store map[string][]Credential // user -> data
}

// NewMemory создаёт storage
func NewMemory() *MemoryStorage {
	return &MemoryStorage{
		store: make(map[string][]Credential),
	}
}

// Save сохраняет данные
func (m *MemoryStorage) Save(user string, c Credential) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store[user] = append(m.store[user], c)
}

// List возвращает данные пользователя
func (m *MemoryStorage) List(user string) []Credential {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.store[user]
}
