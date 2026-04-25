package data

import (
	"encoding/json"
	"net/http"

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

// CreateCredentialRequest запрос
type CreateCredentialRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Meta     string `json:"meta"`
}

// Create создаёт credential
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.UserKey).(string)

	var req CreateCredentialRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	c := storage.Credential{
		ID:       uuid.NewString(),
		User:     user,
		Login:    req.Login,
		Password: req.Password,
		Meta:     req.Meta,
	}

	h.storage.Save(user, c)

	w.WriteHeader(http.StatusCreated)
}

// List возвращает все данные пользователя
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.UserKey).(string)

	data := h.storage.List(user)

	_ = json.NewEncoder(w).Encode(data)
}
