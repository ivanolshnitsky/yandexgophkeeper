package auth

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// AuthService интерфейс auth сервиса.
type AuthService interface {
	Register(username, password string) error
	Login(username, password string) error
}

// User интерфейс
type UserRepo interface {
	Create(User) error
	GetByUsername(string) (User, error)
}

// Service бизнес-логика auth.
type Service struct {
	repo UserRepo
}

// NewService создаёт сервис.
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Register создаёт пользователя.
func (s *Service) Register(username, password string) error {

	_, err := s.repo.GetByUsername(username)
	if err == nil {
		return errors.New("user already exists")
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	return s.repo.Create(User{
		Username:     username,
		PasswordHash: string(hash),
	})
}

// Login проверяет пароль.
func (s *Service) Login(username, password string) error {

	u, err := s.repo.GetByUsername(username)
	if err != nil {
		return errors.New("user not found")
	}

	return bcrypt.CompareHashAndPassword(
		[]byte(u.PasswordHash),
		[]byte(password),
	)
}
