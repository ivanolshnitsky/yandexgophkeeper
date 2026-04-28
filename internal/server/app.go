package server

import (
	"net/http"
	"yandexgophkeeper/internal/config"
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

func New(log *zap.Logger, authHandler *auth.Handler, cfg *config.AppConfig) *App {
	r := chi.NewRouter()

	var store storage.Storage

	if cfg.StorageType == "postgres" {
		pg, err := storage.NewPostgres(cfg.DatabaseDSN)
		if err != nil {
			log.Fatal("postgres init failed", zap.Error(err))
		}
		store = pg
	} else {
		store = storage.NewMemory()
	}

	dataHandler := data.NewHandler(store)

	app := &App{
		router:      r,
		log:         log,
		authHandler: authHandler,
		dataHandler: dataHandler,
		secret:      cfg.SecretKey,
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

		r.Route("/data", func(r chi.Router) {
			r.Post("/", a.dataHandler.Create)
			r.Get("/", a.dataHandler.List)
			r.Get("/{id}", a.dataHandler.Get)
			r.Put("/{id}", a.dataHandler.Update)
			r.Delete("/{id}", a.dataHandler.Delete)
		})
	})
}
