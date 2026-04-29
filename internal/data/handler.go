package data

import (
	"encoding/json"
	"net/http"
	"time"

	"yandexgophkeeper/internal/auth"
	"yandexgophkeeper/internal/crypto"
	"yandexgophkeeper/internal/storage"

	"github.com/google/uuid"
)

// Handler обрабатывает данные пользователя.
type Handler struct {
	storage storage.Storage
	crypto  *crypto.Service
}

// NewHandler создаёт handler.
func NewHandler(s storage.Storage, c *crypto.Service) *Handler {
	return &Handler{s, c}
}

// CreateRequest универсальный запрос.
type CreateRequest struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
	Meta  string          `json:"meta"`
}

type DataResponse struct {
	ID    string          `json:"id"`
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
	Meta  string          `json:"meta"`
}

func user(r *http.Request) string {
	v := r.Context().Value(auth.UserKey)
	if v == nil {
		return ""
	}
	u, _ := v.(string)
	return u
}

// Create создаёт новую запись пользователя.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user := user(r)
	if user == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var plain string
	if err := json.Unmarshal(req.Value, &plain); err != nil {
		http.Error(w, "invalid value", http.StatusBadRequest)
		return
	}

	enc, err := h.crypto.Encrypt([]byte(plain))
	if err != nil {
		http.Error(w, "encrypt error", http.StatusInternalServerError)
		return
	}

	val, _ := json.Marshal(enc)

	d := storage.Data{
		ID:        uuid.NewString(),
		User:      user,
		Type:      storage.DataType(req.Type),
		Value:     val,
		Meta:      req.Meta,
		UpdatedAt: time.Now(),
	}

	if err := h.storage.Save(user, d); err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// List возвращает все записи пользователя.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	user := user(r)

	data, err := h.storage.List(user)
	if err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	resp := make([]DataResponse, 0, len(data))

	for _, d := range data {
		resp = append(resp, DataResponse{
			ID:    d.ID,
			Type:  string(d.Type),
			Value: d.Value,
			Meta:  d.Meta,
		})
	}

	_ = json.NewEncoder(w).Encode(resp)
}

// Get возвращает запись по ID.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	user := user(r)
	id := r.URL.Query().Get("id")

	d, ok, _ := h.storage.GetByID(user, id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	_ = json.NewEncoder(w).Encode(DataResponse{
		ID:    d.ID,
		Type:  string(d.Type),
		Value: d.Value,
		Meta:  d.Meta,
	})
}

// Update обновляет запись пользователя.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	user := user(r)
	id := r.URL.Query().Get("id")

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var plain string
	if err := json.Unmarshal(req.Value, &plain); err != nil {
		http.Error(w, "invalid value", http.StatusBadRequest)
		return
	}

	enc, err := h.crypto.Encrypt([]byte(plain))
	if err != nil {
		http.Error(w, "encrypt error", http.StatusInternalServerError)
		return
	}

	val, _ := json.Marshal(enc)

	d := storage.Data{
		ID:        id,
		User:      user,
		Type:      storage.DataType(req.Type),
		Value:     val,
		Meta:      req.Meta,
		UpdatedAt: time.Now(),
	}

	ok, err := h.storage.Update(user, d)
	if err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete удаляет запись пользователя.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	user := user(r)
	id := r.URL.Query().Get("id")

	ok, err := h.storage.Delete(user, id)
	if err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
