package auth

import (
	"encoding/json"
	"net/http"
)

/*
Handler обрабатывает HTTP-запросы аутентификации.

Отвечает за:
- регистрацию пользователя
- логин пользователя
*/
type Handler struct {
	service *Service
	jwt     *JWTService
}

// NewHandler создаёт новый Handler.
func NewHandler(service *Service, jwt *JWTService) *Handler {
	return &Handler{
		service: service,
		jwt:     jwt,
	}
}

// RegisterRequest содержит данные для регистрации и логина.
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

/*
Register регистрирует нового пользователя.

Request:
POST /register

	{
	  "username": "user",
	  "password": "pass"
	}

Response:
201 Created
*/
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var in RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.Register(in.Username, in.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

/*
Login аутентифицирует пользователя и возвращает JWT-токен.

Request:
POST /login

	{
	  "username": "user",
	  "password": "pass"
	}

Response:
200 OK

	{
	  "token": "jwt..."
	}
*/
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var in RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.Login(in.Username, in.Password); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := h.jwt.Generate(in.Username)
	if err != nil {
		http.Error(w, "token error", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
