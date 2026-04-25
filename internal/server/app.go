package server

import (
	"net/http"

	"yandexgophkeeper/internal/logger"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type App struct {
	router *chi.Mux
	log    *zap.Logger
}

func New(log *zap.Logger) *App {
	r := chi.NewRouter()

	app := &App{
		router: r,
		log:    log,
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
}
