package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
)

// PostgresStorage реализует Storage через PostgreSQL.
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgres создаёт подключение к БД и запускает миграции.
func NewPostgres(dsn string) (*PostgresStorage, error) {
	if err := runMigrations(dsn); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

// runMigrations применяет миграции.
func runMigrations(dsn string) error {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}

// Save сохраняет данные пользователя.
func (p *PostgresStorage) Save(user string, d Data) error {
	_, err := p.db.Exec(
		`INSERT INTO data (id, user_id, type, value, meta)
		 VALUES ($1,$2,$3,$4,$5)`,
		d.ID, user, d.Type, d.Value, d.Meta,
	)
	if err != nil {
		return fmt.Errorf("insert data: %w", err)
	}

	return nil
}

// List возвращает все данные пользователя.
func (p *PostgresStorage) List(user string) ([]Data, error) {
	rows, err := p.db.Query(
		`SELECT id, user_id, type, value, meta FROM data WHERE user_id=$1`,
		user,
	)
	if err != nil {
		return nil, fmt.Errorf("query list: %w", err)
	}
	defer rows.Close()

	var res []Data

	for rows.Next() {
		var d Data
		if err := rows.Scan(&d.ID, &d.User, &d.Type, &d.Value, &d.Meta); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		res = append(res, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return res, nil
}

// GetByID возвращает запись по ID.
func (p *PostgresStorage) GetByID(user, id string) (Data, bool, error) {
	var d Data

	err := p.db.QueryRow(
		`SELECT id, user_id, type, value, meta 
		 FROM data WHERE user_id=$1 AND id=$2`,
		user, id,
	).Scan(&d.ID, &d.User, &d.Type, &d.Value, &d.Meta)

	if errors.Is(err, sql.ErrNoRows) {
		return Data{}, false, nil
	}
	if err != nil {
		return Data{}, false, fmt.Errorf("get by id: %w", err)
	}

	return d, true, nil
}

// Update обновляет запись.
func (p *PostgresStorage) Update(user string, d Data) (bool, error) {
	res, err := p.db.Exec(
		`UPDATE data 
		 SET type=$1, value=$2, meta=$3 
		 WHERE id=$4 AND user_id=$5`,
		d.Type, d.Value, d.Meta, d.ID, user,
	)
	if err != nil {
		return false, fmt.Errorf("update data: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("rows affected: %w", err)
	}

	return rows > 0, nil
}

// Delete удаляет запись.
func (p *PostgresStorage) Delete(user, id string) (bool, error) {
	res, err := p.db.Exec(
		`DELETE FROM data WHERE id=$1 AND user_id=$2`,
		id, user,
	)
	if err != nil {
		return false, fmt.Errorf("delete data: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("rows affected: %w", err)
	}

	return rows > 0, nil
}
