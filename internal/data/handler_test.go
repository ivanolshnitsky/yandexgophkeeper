package data

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"yandexgophkeeper/internal/auth"
	"yandexgophkeeper/internal/crypto"
	"yandexgophkeeper/internal/storage"
)

type testStorage struct {
	items map[string]storage.Data
}

func newTestStorage() *testStorage {
	return &testStorage{items: map[string]storage.Data{}}
}

func (s *testStorage) Save(user string, d storage.Data) error {
	s.items[d.ID] = d
	return nil
}
func (s *testStorage) List(user string) ([]storage.Data, error) {
	res := []storage.Data{}
	for _, v := range s.items {
		res = append(res, v)
	}
	return res, nil
}
func (s *testStorage) GetByID(user, id string) (storage.Data, bool, error) {
	d, ok := s.items[id]
	return d, ok, nil
}
func (s *testStorage) Update(user string, d storage.Data) (bool, error) {
	_, ok := s.items[d.ID]
	if !ok {
		return false, nil
	}
	s.items[d.ID] = d
	return true, nil
}
func (s *testStorage) Delete(user, id string) (bool, error) {
	_, ok := s.items[id]
	if !ok {
		return false, nil
	}
	delete(s.items, id)
	return true, nil
}
func (s *testStorage) Close() error { return nil }

func setup() *Handler {
	return NewHandler(newTestStorage(), crypto.New("secret"))
}

func withUser(r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), auth.UserKey, "user"))
}

func createItem(t *testing.T, h *Handler) string {
	body := `{"type":"text","value":"hello","meta":"m"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req = withUser(req)

	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("create failed: %s", w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req = withUser(req)
	w = httptest.NewRecorder()
	h.List(w, req)

	var list []DataResponse
	_ = json.Unmarshal(w.Body.Bytes(), &list)

	return list[0].ID
}

func TestCreateAndList_Decrypted(t *testing.T) {
	h := setup()

	id := createItem(t, h)
	if id == "" {
		t.Fatal("id empty")
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = withUser(req)

	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatal("list status not ok")
	}

	var list []DataResponse
	_ = json.Unmarshal(w.Body.Bytes(), &list)

	var value string
	_ = json.Unmarshal(list[0].Value, &value)

	if value != "hello" {
		t.Fatal("value not decrypted")
	}
}

func TestGet(t *testing.T) {
	h := setup()
	id := createItem(t, h)

	req := httptest.NewRequest(http.MethodGet, "/"+id, nil)
	req = withUser(req)
	req = muxID(req, id)

	w := httptest.NewRecorder()
	h.Get(w, req)

	if w.Code != http.StatusOK {
		t.Fatal("get failed")
	}
}

func TestUpdate(t *testing.T) {
	h := setup()
	id := createItem(t, h)

	body := `{"type":"text","value":"updated","meta":"m"}`
	req := httptest.NewRequest(http.MethodPut, "/"+id, bytes.NewBufferString(body))
	req = withUser(req)
	req = muxID(req, id)

	w := httptest.NewRecorder()
	h.Update(w, req)

	if w.Code != http.StatusOK {
		t.Fatal("update failed")
	}
}

func TestDelete(t *testing.T) {
	h := setup()
	id := createItem(t, h)

	req := httptest.NewRequest(http.MethodDelete, "/"+id, nil)
	req = withUser(req)
	req = muxID(req, id)

	w := httptest.NewRecorder()
	h.Delete(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatal("delete failed")
	}
}

func muxID(r *http.Request, id string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
