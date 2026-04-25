package auth

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockService struct {
	registerErr error
	loginErr    error
}

func (m *mockService) Register(username, password string) error {
	return m.registerErr
}

func (m *mockService) Login(username, password string) error {
	return m.loginErr
}

func setupHandler() *Handler {
	return NewHandler(
		&mockService{},
		NewJWTService("secret", time.Hour),
	)
}

func TestRegister(t *testing.T) {
	h := setupHandler()

	body := []byte(`{"username":"test","password":"123"}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/register",
		bytes.NewBuffer(body),
	)

	w := httptest.NewRecorder()

	h.Register(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRegister_Empty(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest(
		http.MethodPost,
		"/register",
		bytes.NewBuffer([]byte(`{}`)),
	)

	w := httptest.NewRecorder()

	h.Register(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin(t *testing.T) {
	h := setupHandler()

	body := []byte(`{"username":"test","password":"123"}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		bytes.NewBuffer(body),
	)

	w := httptest.NewRecorder()

	h.Login(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

func TestLogin_Invalid(t *testing.T) {
	h := NewHandler(
		&mockService{
			loginErr: errors.New("invalid"),
		},
		NewJWTService("secret", time.Hour),
	)

	body := []byte(`{"username":"test","password":"wrong"}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		bytes.NewBuffer(body),
	)

	w := httptest.NewRecorder()

	h.Login(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
