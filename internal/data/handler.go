package data

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"yandexgophkeeper/internal/auth"
	"yandexgophkeeper/internal/crypto"
	"yandexgophkeeper/internal/storage"
)

// Handler обрабатывает данные пользователя.
type Handler struct {
	storage storage.Storage
	crypto  *crypto.Service
}

// NewHandler создаёт handler.
func NewHandler(s storage.Storage, c *crypto.Service) *Handler {
	return &Handler{
		storage: s,
		crypto:  c,
	}
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

func userFromCtx(r *http.Request) string {
	v := r.Context().Value(auth.UserKey)
	if v == nil {
		return ""
	}
	u, _ := v.(string)
	return u
}

// Create создаёт новую запись пользователя.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	u := userFromCtx(r)
	if u == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	enc, err := h.crypto.Encrypt(req.Value)
	if err != nil {
		http.Error(w, "encrypt error", http.StatusInternalServerError)
		return
	}

	encJSON, err := json.Marshal(enc)
	if err != nil {
		http.Error(w, "marshal error", http.StatusInternalServerError)
		return
	}

	d := storage.Data{
		ID:        uuid.NewString(),
		User:      u,
		Type:      storage.DataType(req.Type),
		Value:     encJSON,
		Meta:      req.Meta,
		UpdatedAt: time.Now(),
	}

	if err := h.storage.Save(u, d); err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// List возвращает все записи пользователя.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	u := userFromCtx(r)
	if u == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	items, err := h.storage.List(u)
	if err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	resp := make([]DataResponse, 0, len(items))

	for _, item := range items {

		var enc string
		if err := json.Unmarshal(item.Value, &enc); err != nil {
			http.Error(w, "decode error", http.StatusInternalServerError)
			return
		}

		dec, err := h.crypto.Decrypt(enc)
		if err != nil {
			http.Error(w, "decrypt error", http.StatusInternalServerError)
			return
		}

		resp = append(resp, DataResponse{
			ID:    item.ID,
			Type:  string(item.Type),
			Value: json.RawMessage(dec),
			Meta:  item.Meta,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// Get возвращает запись по ID.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	u := userFromCtx(r)
	if u == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	item, ok, err := h.storage.GetByID(u, id)
	if err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.NotFound(w, r)
		return
	}

	var enc string
	if err := json.Unmarshal(item.Value, &enc); err != nil {
		http.Error(w, "decode error", http.StatusInternalServerError)
		return
	}

	dec, err := h.crypto.Decrypt(enc)
	if err != nil {
		http.Error(w, "decrypt error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(DataResponse{
		ID:    item.ID,
		Type:  string(item.Type),
		Value: json.RawMessage(dec),
		Meta:  item.Meta,
	})
}

// Update обновляет запись пользователя.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	u := userFromCtx(r)
	if u == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	enc, err := h.crypto.Encrypt(req.Value)
	if err != nil {
		http.Error(w, "encrypt error", http.StatusInternalServerError)
		return
	}

	encJSON, err := json.Marshal(enc)
	if err != nil {
		http.Error(w, "marshal error", http.StatusInternalServerError)
		return
	}

	d := storage.Data{
		ID:        id,
		User:      u,
		Type:      storage.DataType(req.Type),
		Value:     encJSON,
		Meta:      req.Meta,
		UpdatedAt: time.Now(),
	}

	ok, err := h.storage.Update(u, d)
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
	u := userFromCtx(r)
	id := chi.URLParam(r, "id")

	ok, err := h.storage.Delete(u, id)
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
