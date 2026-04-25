package auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupHandler() *Handler {
	service := NewService()
	jwt := NewJWTService("secret")
	return NewHandler(service, jwt)
}

func TestRegister(t *testing.T) {
	h := setupHandler()

	body := []byte(`{"username":"test","password":"123"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Register(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestLogin(t *testing.T) {
	h := setupHandler()

	_ = h.service.Register("test", "123")

	body := []byte(`{"username":"test","password":"123"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Login(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}
