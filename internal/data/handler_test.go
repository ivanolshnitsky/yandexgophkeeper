package data

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"yandexgophkeeper/internal/auth"
	"yandexgophkeeper/internal/crypto"
	"yandexgophkeeper/internal/storage"
)

func setup() (*Handler, storage.Storage) {
	st := storage.NewMemory()
	cr := crypto.New("test-secret-key-1234567890")
	return NewHandler(st, cr), st
}

func withUser(r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), auth.UserKey, "user"))
}

func createItem(t *testing.T, h *Handler) string {
	t.Helper()

	body := `{"type":"text","value":"\"hello\"","meta":"m"}`

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req = withUser(req)

	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("create failed: %d body=%s", w.Code, w.Body.String())
	}

	// list items
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req = withUser(req)

	w = httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("list failed: %d body=%s", w.Code, w.Body.String())
	}

	var list []DataResponse
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatalf("decode list failed: %v body=%s", err, w.Body.String())
	}

	if len(list) == 0 {
		t.Fatal("empty list after create")
	}

	return list[0].ID
}

func TestCreate(t *testing.T) {
	h, _ := setup()

	body := `{"type":"text","value":"\"hello\"","meta":"m"}`

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req = withUser(req)

	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestList(t *testing.T) {
	h, _ := setup()

	_ = createItem(t, h)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = withUser(req)

	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	var list []DataResponse
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatalf("json error: %v body=%s", err, w.Body.String())
	}

	if len(list) == 0 {
		t.Fatal("expected at least 1 item")
	}
}

func TestGet(t *testing.T) {
	h, _ := setup()

	id := createItem(t, h)

	req := httptest.NewRequest(http.MethodGet, "/?id="+id, nil)
	req = withUser(req)

	w := httptest.NewRecorder()
	h.Get(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	var resp DataResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json error: %v body=%s", err, w.Body.String())
	}

	if resp.ID != id {
		t.Fatalf("wrong id: got %s want %s", resp.ID, id)
	}
}

func TestUpdate(t *testing.T) {
	h, _ := setup()

	id := createItem(t, h)

	body := `{"type":"text","value":"\"updated\"","meta":"m2"}`

	req := httptest.NewRequest(http.MethodPut, "/?id="+id, bytes.NewBufferString(body))
	req = withUser(req)

	w := httptest.NewRecorder()
	h.Update(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestDelete(t *testing.T) {
	h, _ := setup()

	id := createItem(t, h)

	req := httptest.NewRequest(http.MethodDelete, "/?id="+id, nil)
	req = withUser(req)

	w := httptest.NewRecorder()
	h.Delete(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 got %d body=%s", w.Code, w.Body.String())
	}
}
