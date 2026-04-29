package server

import (
	"net/http"

	"yandexgophkeeper/internal/auth"
	"yandexgophkeeper/internal/config"
	"yandexgophkeeper/internal/crypto"
	"yandexgophkeeper/internal/data"
	"yandexgophkeeper/internal/logger"
	"yandexgophkeeper/internal/storage"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type App struct {
	router      *chi.Mux
	log         *zap.Logger
	authHandler *auth.Handler
	dataHandler *data.Handler
	secret      string
	store       storage.Storage
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

	cryptoService := crypto.New(cfg.SecretKey)
	dataHandler := data.NewHandler(store, cryptoService)

	app := &App{
		router:      r,
		log:         log,
		authHandler: authHandler,
		dataHandler: dataHandler,
		secret:      cfg.SecretKey,
		store:       store,
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
	a.router.Post("/register", a.authHandler.Register)
	a.router.Post("/login", a.authHandler.Login)
	a.router.Get("/ping", a.handlePing)

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

// Close аккуратно закрывает storage если он поддерживает Close
func (a *App) Close() error {
	if a.store == nil {
		return nil
	}

	if closer, ok := a.store.(interface{ Close() error }); ok {
		return closer.Close()
	}

	return nil
}
