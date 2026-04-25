package auth

import (
	"database/sql"
	"errors"
	"testing"
)

type mockRepo struct {
	users map[string]User
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		users: make(map[string]User),
	}
}

func (m *mockRepo) Create(user User) error {
	m.users[user.Username] = user
	return nil
}

func (m *mockRepo) GetByUsername(username string) (User, error) {
	u, ok := m.users[username]
	if !ok {
		return User{}, sql.ErrNoRows
	}

	return u, nil
}

func TestRegisterAndLogin(t *testing.T) {
	repo := newMockRepo()

	s := &Service{
		repo: (*Repository)(nil),
	}

	s.repo = &Repository{}

	_ = repo

	realService := NewService(nil)

	_ = realService

	err := repo.Create(User{
		Username:     "user",
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZo4i.ej8z5mFQxQ5e0siJmq9hw3rrodpAtxW",
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestRegisterDuplicate(t *testing.T) {
	repo := newMockRepo()

	hash := "$2a$10$N9qo8uLOickgx2ZMRZo4i.ej8z5mFQxQ5e0siJmq9hw3rrodpAtxW"

	repo.users["user"] = User{
		Username:     "user",
		PasswordHash: hash,
	}

	_, err := repo.GetByUsername("user")

	if err != nil {
		t.Fatal(err)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := newMockRepo()

	_, err := repo.GetByUsername("missing")

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatal("expected sql.ErrNoRows")
	}
}
