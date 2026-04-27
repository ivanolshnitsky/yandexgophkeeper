package data

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"yandexgophkeeper/internal/auth"
	"yandexgophkeeper/internal/storage"

	"github.com/google/uuid"
)

// Handler обрабатывает данные пользователя
type Handler struct {
	storage *storage.MemoryStorage
}

// NewHandler создаёт handler
func NewHandler(s *storage.MemoryStorage) *Handler {
	return &Handler{storage: s}
}

// CreateRequest универсальный запрос
type CreateRequest struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
	Meta  string          `json:"meta"`
}

// Create создаёт новую запись пользователя.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.UserKey).(string)

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Type == "" || len(req.Value) == 0 {
		http.Error(w, "type and value required", http.StatusBadRequest)
		return
	}

	d := storage.Data{
		ID:    uuid.NewString(),
		User:  user,
		Type:  storage.DataType(req.Type),
		Value: req.Value,
		Meta:  req.Meta,
	}

	h.storage.Save(user, d)

	w.WriteHeader(http.StatusCreated)
}

// List возвращает все записи пользователя.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.UserKey).(string)

	data := h.storage.List(user)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

// Get возвращает запись по ID.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.UserKey).(string)
	id := chi.URLParam(r, "id")

	d, ok := h.storage.GetByID(user, id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	_ = json.NewEncoder(w).Encode(d)
}

// Update обновляет запись пользователя.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.UserKey).(string)
	id := chi.URLParam(r, "id")

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Type == "" || len(req.Value) == 0 {
		http.Error(w, "type and value required", http.StatusBadRequest)
		return
	}

	d := storage.Data{
		ID:    id,
		User:  user,
		Type:  storage.DataType(req.Type),
		Value: req.Value,
		Meta:  req.Meta,
	}

	if !h.storage.Update(user, d) {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete удаляет запись пользователя.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.UserKey).(string)
	id := chi.URLParam(r, "id")

	if !h.storage.Delete(user, id) {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
