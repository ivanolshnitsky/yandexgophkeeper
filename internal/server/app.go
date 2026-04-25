package server

import (
	"net/http"
	"yandexgophkeeper/internal/data"
	"yandexgophkeeper/internal/storage"

	"yandexgophkeeper/internal/auth"
	"yandexgophkeeper/internal/logger"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type App struct {
	router *chi.Mux
	log    *zap.Logger

	authHandler *auth.Handler
	dataHandler *data.Handler

	secret string
}

func New(log *zap.Logger, authHandler *auth.Handler, secret string) *App {
	r := chi.NewRouter()

	store := storage.NewMemory()
	dataHandler := data.NewHandler(store)

	app := &App{
		router:      r,
		log:         log,
		authHandler: authHandler,
		dataHandler: dataHandler,
		secret:      secret,
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
	// public
	a.router.Post("/register", a.authHandler.Register)
	a.router.Post("/login", a.authHandler.Login)
	a.router.Get("/ping", a.handlePing)

	// protected
	a.router.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return auth.Middleware(a.secret, next)
		})

		r.Post("/data", a.dataHandler.Create)
		r.Get("/data", a.dataHandler.List)
	})
}
