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

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgres(dsn string) (*PostgresStorage, error) {
	if err := runMigrations(dsn); err != nil {
		return nil, fmt.Errorf("migrations: %w", err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}

func runMigrations(dsn string) error {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (p *PostgresStorage) Save(user string, d Data) error {
	_, err := p.db.Exec(
		`INSERT INTO data (id, user_id, type, value, meta, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		d.ID, user, d.Type, d.Value, d.Meta, d.UpdatedAt,
	)
	return err
}

func (p *PostgresStorage) List(user string) ([]Data, error) {
	rows, err := p.db.Query(
		`SELECT id, user_id, type, value, meta, updated_at 
		 FROM data WHERE user_id=$1`,
		user,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Data

	for rows.Next() {
		var d Data
		if err := rows.Scan(
			&d.ID,
			&d.User,
			&d.Type,
			&d.Value,
			&d.Meta,
			&d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, d)
	}

	return res, rows.Err()
}

func (p *PostgresStorage) GetByID(user, id string) (Data, bool, error) {
	var d Data

	err := p.db.QueryRow(
		`SELECT id, user_id, type, value, meta, updated_at
		 FROM data WHERE user_id=$1 AND id=$2`,
		user, id,
	).Scan(&d.ID, &d.User, &d.Type, &d.Value, &d.Meta, &d.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return Data{}, false, nil
	}

	if err != nil {
		return Data{}, false, err
	}

	return d, true, nil
}

// Update обновляет запись.
func (p *PostgresStorage) Update(user string, d Data) (bool, error) {
	res, err := p.db.Exec(
		`UPDATE data 
		 SET type=$1, value=$2, meta=$3, updated_at=$4
		 WHERE id=$5 AND user_id=$6`,
		d.Type, d.Value, d.Meta, d.UpdatedAt, d.ID, user,
	)
	if err != nil {
		return false, fmt.Errorf("update data: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("rows affected: %w", err)
	}

	return n > 0, nil
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

	n, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("rows affected: %w", err)
	}

	return n > 0, nil
}

// Close закрывает соединение.
func (p *PostgresStorage) Close() error {
	return p.db.Close()
}
