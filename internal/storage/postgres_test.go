package storage

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func newMockStorage(t *testing.T) (*PostgresStorage, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}

	st := &PostgresStorage{db: db}
	return st, mock, func() { db.Close() }
}

func testData() Data {
	return Data{
		ID:        "id1",
		User:      "u1",
		Type:      "text",
		Value:     []byte("hello"),
		Meta:      "{}",
		UpdatedAt: time.Now(),
	}
}

func TestPostgres_Save(t *testing.T) {
	st, mock, closeFn := newMockStorage(t)
	defer closeFn()

	d := testData()

	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO data (id, user_id, type, value, meta, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
	)).
		WithArgs(d.ID, "u1", d.Type, d.Value, d.Meta, d.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := st.Save("u1", d); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPostgres_List(t *testing.T) {
	st, mock, closeFn := newMockStorage(t)
	defer closeFn()

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "type", "value", "meta", "updated_at",
	}).AddRow("id1", "u1", "text", []byte("hello"), "{}", time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, user_id, type, value, meta, updated_at 
		 FROM data WHERE user_id=$1`,
	)).
		WithArgs("u1").
		WillReturnRows(rows)

	res, err := st.List("u1")
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 row got %d", len(res))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPostgres_GetByID_Found(t *testing.T) {
	st, mock, closeFn := newMockStorage(t)
	defer closeFn()

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "type", "value", "meta", "updated_at",
	}).AddRow("id1", "u1", "text", []byte("hello"), "{}", time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, user_id, type, value, meta, updated_at
		 FROM data WHERE user_id=$1 AND id=$2`,
	)).
		WithArgs("u1", "id1").
		WillReturnRows(rows)

	_, ok, err := st.GetByID("u1", "id1")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected found")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPostgres_GetByID_NotFound(t *testing.T) {
	st, mock, closeFn := newMockStorage(t)
	defer closeFn()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, user_id, type, value, meta, updated_at
		 FROM data WHERE user_id=$1 AND id=$2`,
	)).
		WithArgs("u1", "missing").
		WillReturnError(sql.ErrNoRows)

	_, ok, err := st.GetByID("u1", "missing")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("should be not found")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPostgres_Update(t *testing.T) {
	st, mock, closeFn := newMockStorage(t)
	defer closeFn()

	d := testData()

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE data 
		 SET type=$1, value=$2, meta=$3, updated_at=$4
		 WHERE id=$5 AND user_id=$6`,
	)).
		WithArgs(d.Type, d.Value, d.Meta, d.UpdatedAt, d.ID, "u1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	ok, err := st.Update("u1", d)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected updated")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPostgres_Delete(t *testing.T) {
	st, mock, closeFn := newMockStorage(t)
	defer closeFn()

	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM data WHERE id=$1 AND user_id=$2`,
	)).
		WithArgs("id1", "u1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	ok, err := st.Delete("u1", "id1")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected deleted")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
