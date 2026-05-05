package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"yandexgophkeeper/internal/auth"
	"yandexgophkeeper/internal/config"
	"yandexgophkeeper/internal/storage"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type testStorage struct{}

func (s *testStorage) Save(user string, d storage.Data) error   { return nil }
func (s *testStorage) List(user string) ([]storage.Data, error) { return nil, nil }
func (s *testStorage) GetByID(user, id string) (storage.Data, bool, error) {
	return storage.Data{}, false, nil
}
func (s *testStorage) Update(user string, d storage.Data) (bool, error) { return true, nil }
func (s *testStorage) Delete(user, id string) (bool, error)             { return true, nil }
func (s *testStorage) Close() error                                     { return nil }

func TestHandlePing(t *testing.T) {
	log := zap.NewNop()

	cfg := &config.AppConfig{
		JWTSecret: "test-secret",
		CryptoKey: "crypto-secret",
		TokenTTL:  time.Hour,
	}

	authHandler := auth.NewHandler(nil, nil)

	app := New(log, authHandler, cfg, &testStorage{})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	app.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}
