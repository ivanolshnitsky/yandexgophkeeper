package auth

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO users(username, password_hash)
		 VALUES ($1, $2)`,
	)).
		WithArgs("user", "hash").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(User{
		Username:     "user",
		PasswordHash: "hash",
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestRepository_GetByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{
		"username",
		"password_hash",
	}).AddRow("user", "hash")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT username, password_hash
		 FROM users
		 WHERE username=$1`,
	)).
		WithArgs("user").
		WillReturnRows(rows)

	user, err := repo.GetByUsername("user")
	if err != nil {
		t.Fatal(err)
	}

	if user.Username != "user" {
		t.Fatalf("expected user got %s", user.Username)
	}

	if user.PasswordHash != "hash" {
		t.Fatalf("expected hash got %s", user.PasswordHash)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestRepository_GetByUsername_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT username, password_hash
		 FROM users
		 WHERE username=$1`,
	)).
		WithArgs("missing").
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetByUsername("missing")
	if err == nil {
		t.Fatal("expected error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
