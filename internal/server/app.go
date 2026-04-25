package server

import (
	"net/http"

	"yandexgophkeeper/internal/auth"
	"yandexgophkeeper/internal/logger"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type App struct {
	router *chi.Mux
	log    *zap.Logger

	authHandler *auth.Handler
}

func New(log *zap.Logger, authHandler *auth.Handler) *App {
	r := chi.NewRouter()

	app := &App{
		router:      r,
		log:         log,
		authHandler: authHandler,
	}

	r.Use(func(next http.Handler) http.Handler {
		return logger.Middleware(log, next)
	})

	app.routes()

	return app
}

func (a *App) Router() http.Handler {
	return a.router
}

func (a *App) routes() {
	a.router.Get("/ping", a.handlePing)

	a.router.Post("/register", a.authHandler.Register)
	a.router.Post("/login", a.authHandler.Login)
}
