package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"yandexgophkeeper/internal/auth"
	"yandexgophkeeper/internal/config"
	"yandexgophkeeper/internal/logger"
	"yandexgophkeeper/internal/server"
	"yandexgophkeeper/internal/storage"

	"go.uber.org/zap"
)

func main() {

	cfg, err := config.New()
	if err != nil {
		panic("config error: " + err.Error())
	}

	log, err := logger.Initialize(cfg.LogLevel)
	if err != nil {
		panic("logger init error: " + err.Error())
	}
	defer func() {
		_ = log.Sync()
	}()

	store, err := storage.NewPostgres(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal("postgres init failed", zap.Error(err))
	}

	repo := auth.NewRepository(store.DB())
	authService := auth.NewService(repo)
	jwtService := auth.NewJWTService(cfg.JWTSecret, cfg.TokenTTL)

	authHandler := auth.NewHandler(authService, jwtService)

	app := server.New(log, authHandler, cfg, store)

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: app.Router(),
	}

	go func() {
		log.Info("server started",
			zap.String("addr", cfg.ServerAddress),
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = srv.Shutdown(ctx)
	_ = app.Close()

	log.Info("server exited")
}
