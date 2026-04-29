package storage

import "sync"

// MemoryStorage in-memory хранилище
type MemoryStorage struct {
	mu   sync.RWMutex
	data map[string][]Data
}

// NewMemory создаёт storage
func NewMemory() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string][]Data),
	}
}

// Save сохраняет данные
func (s *MemoryStorage) Save(user string, d Data) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[user] = append(s.data[user], d)
	return nil
}

// List возвращает данные пользователя
func (s *MemoryStorage) List(user string) ([]Data, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.data[user], nil
}

// GetByID возвращает запись по id
func (s *MemoryStorage) GetByID(user, id string) (Data, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, d := range s.data[user] {
		if d.ID == id {
			return d, true, nil
		}
	}
	return Data{}, false, nil
}

// Update обновляет запись
func (s *MemoryStorage) Update(user string, d Data) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, v := range s.data[user] {
		if v.ID == d.ID {
			s.data[user][i] = d
			return true, nil
		}
	}
	return false, nil
}

// Delete удаляет запись
func (s *MemoryStorage) Delete(user, id string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	items := s.data[user]
	for i, v := range items {
		if v.ID == id {
			s.data[user] = append(items[:i], items[i+1:]...)
			return true, nil
		}
	}
	return false, nil
}

func (s *MemoryStorage) Close() error {
	return nil
}
