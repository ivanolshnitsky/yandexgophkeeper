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

	"go.uber.org/zap"
)

func main() {
	cfg := config.New()

	log, err := logger.Initialize(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Sync() }()

	authService := auth.NewService()
	jwtService := auth.NewJWTService(cfg.SecretKey)

	authHandler := auth.NewHandler(authService, jwtService)

	app := server.New(log, authHandler, cfg)

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: app.Router(),
	}

	go func() {
		log.Info("server started", zap.String("addr", cfg.ServerAddress))

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("shutdown failed", zap.Error(err))
	}

	log.Info("server exited")
}
