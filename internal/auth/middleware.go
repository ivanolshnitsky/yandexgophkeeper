package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

// UserKey ключ для хранения пользователя в контексте
const UserKey contextKey = "user"

// Middleware проверяет JWT и кладёт username в context
func Middleware(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "no token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid claims", http.StatusUnauthorized)
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			http.Error(w, "invalid username", http.StatusUnauthorized)
			return
		}

		// 🔥 ВОТ ГДЕ ИСПОЛЬЗУЕТСЯ UserKey
		ctx := context.WithValue(r.Context(), UserKey, username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
