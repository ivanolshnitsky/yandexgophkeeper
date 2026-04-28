package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"yandexgophkeeper/internal/auth"
	"yandexgophkeeper/internal/config"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

func TestHandlePing(t *testing.T) {
	log := zap.NewNop()
	cfg := config.New()

	authService := auth.NewService()
	jwtService := auth.NewJWTService("secret")
	authHandler := auth.NewHandler(authService, jwtService)

	app := New(log, authHandler, cfg)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	app.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
