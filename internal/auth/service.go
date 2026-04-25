package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Service бизнес-логика auth
type Service struct {
	users map[string]User
}

// NewService создаёт сервис
func NewService() *Service {
	return &Service{
		users: make(map[string]User),
	}
}

// Register создаёт пользователя
func (s *Service) Register(username, password string) error {
	if _, ok := s.users[username]; ok {
		return errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	s.users[username] = User{
		Username:     username,
		PasswordHash: string(hash),
	}

	return nil
}

// Login проверка пароля
func (s *Service) Login(username, password string) error {
	u, ok := s.users[username]
	if !ok {
		return errors.New("user not found")
	}

	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}
