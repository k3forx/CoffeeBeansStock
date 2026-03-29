package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/api/handlers"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/auth"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/config"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/repository"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/router"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/usecase"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}
}

func parseLogLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	level := parseLogLevel(cfg.LogLevel)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL())
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	slog.Info("database connection established")

	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	userRepo := repository.NewUserRepository(pool)
	coffeeBeanRepo := repository.NewCoffeeBeanRepository(pool)
	usageHistoryRepo := repository.NewUsageHistoryRepository(pool)
	uow := repository.NewUnitOfWork(pool)

	authService := usecase.NewAuthService(userRepo, jwtManager)
	authHandler := handlers.NewAuthHandler(authService)
	coffeeBeansService := usecase.NewCoffeeBeansService(coffeeBeanRepo, uow)
	coffeeBeansHandler := handlers.NewCoffeeBeansHandler(coffeeBeansService)
	usageHistoryService := usecase.NewUsageHistoryService(usageHistoryRepo, coffeeBeanRepo, uow)
	usageHistoryHandler := handlers.NewUsageHistoryHandler(usageHistoryService)

	r := router.New(router.Deps{
		AuthHandler:         authHandler,
		CoffeeBeansHandler:  coffeeBeansHandler,
		UsageHistoryHandler: usageHistoryHandler,
		TokenManager:        jwtManager,
		HealthCheck: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if err := pool.Ping(r.Context()); err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = fmt.Fprintf(w, `{"status":"error","database":"disconnected"}`)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"status":"ok","database":"connected"}`)
		},
	})

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	case <-quit:
	}

	slog.Info("shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}
	slog.Info("server stopped")
	return nil
}
