package data

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"yandexgophkeeper/internal/auth"
	"yandexgophkeeper/internal/storage"
)

func TestList(t *testing.T) {
	st := storage.NewMemory()
	h := NewHandler(st)

	req := httptest.NewRequest(http.MethodGet, "/data", nil)
	req = req.WithContext(
		context.WithValue(req.Context(), auth.UserKey, "user"),
	)

	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatal("expected 200")
	}
}
