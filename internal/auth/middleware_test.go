package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMiddleware_NoToken(t *testing.T) {
	handler := Middleware(
		"secret",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatal("expected unauthorized")
	}
}

func TestMiddleware_WithToken(t *testing.T) {
	j := NewJWTService("secret", time.Hour)

	token, err := j.Generate("user")
	if err != nil {
		t.Fatal(err)
	}

	handler := Middleware(
		"secret",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			user := r.Context().Value(UserKey)
			if user == nil {
				t.Fatal("user not in context")
			}

			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatal("expected 200")
	}
}

func TestMiddleware_InvalidToken(t *testing.T) {
	handler := Middleware(
		"secret",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatal("expected unauthorized")
	}
}
