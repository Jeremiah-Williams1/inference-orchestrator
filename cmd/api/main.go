package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Jeremiah-Williams1/inference-orchestrator/config"
	routes "github.com/Jeremiah-Williams1/inference-orchestrator/internal/router"
	"github.com/Jeremiah-Williams1/inference-orchestrator/pkg/logger"
	"github.com/Jeremiah-Williams1/inference-orchestrator/pkg/redisclient"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env in development. In production, env vars are injected by the platform.
	// _ = means we intentionally ignore the error — .env not existing is fine in production.
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	log := logger.New(cfg.LogLevel, cfg.LogFormat, cfg.Env)
	slog.SetDefault(log)

	redisClient, err := redisclient.New(cfg.RedisURL)
	if err != nil {
		slog.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}
	slog.Info("redis connected")
	defer redisClient.Close()

	// TODO: wire queue and service once implementations are ready
	// redisQueue := queue.NewRedisQueue(redisClient)
	// jobSvc := service.NewJobService(redisQueue)

	// nil is safe for now — health endpoint does not use the service
	srv := routes.New(cfg, log, redisClient)

	httpSrv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           srv.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("server starting", "port", cfg.Port, "env", cfg.Env)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}
