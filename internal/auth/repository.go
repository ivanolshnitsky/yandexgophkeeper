package auth

import "database/sql"

// Repository репозиторий пользователей.
type Repository struct {
	db *sql.DB
}

// NewRepository создаёт Repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create создаёт пользователя.
func (r *Repository) Create(user User) error {
	_, err := r.db.Exec(
		`INSERT INTO users(username, password_hash)
		 VALUES ($1, $2)`,
		user.Username,
		user.PasswordHash,
	)

	return err
}

// GetByUsername возвращает пользователя.
func (r *Repository) GetByUsername(username string) (User, error) {
	var u User

	err := r.db.QueryRow(
		`SELECT username, password_hash
		 FROM users
		 WHERE username=$1`,
		username,
	).Scan(&u.Username, &u.PasswordHash)

	return u, err
}
