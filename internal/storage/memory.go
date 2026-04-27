package storage

import "sync"

// MemoryStorage in-memory хранилище
type MemoryStorage struct {
	mu    sync.RWMutex
	store map[string][]Data
}

// NewMemory создаёт storage
func NewMemory() *MemoryStorage {
	return &MemoryStorage{
		store: make(map[string][]Data),
	}
}

// Save сохраняет данные
func (m *MemoryStorage) Save(user string, d Data) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store[user] = append(m.store[user], d)
}

// List возвращает данные пользователя
func (m *MemoryStorage) List(user string) []Data {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.store[user]
}

// GetByID возвращает запись по id
func (m *MemoryStorage) GetByID(user, id string) (Data, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, d := range m.store[user] {
		if d.ID == id {
			return d, true
		}
	}
	return Data{}, false
}

// Update обновляет запись
func (m *MemoryStorage) Update(user string, updated Data) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	items := m.store[user]

	for i, d := range items {
		if d.ID == updated.ID {
			items[i] = updated
			m.store[user] = items
			return true
		}
	}
	return false
}

// Delete удаляет запись
func (m *MemoryStorage) Delete(user, id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	items := m.store[user]

	for i, d := range items {
		if d.ID == id {
			m.store[user] = append(items[:i], items[i+1:]...)
			return true
		}
	}
	return false
}
