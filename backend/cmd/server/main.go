package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api/handlers"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/auth"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/config"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/database"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/repository"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/router"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/services"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL())
	if err != nil {
		slog.Error("failed to create connection pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("database connection established")

	queries := database.New(pool)
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	userRepo := repository.NewUserRepository(queries)
	coffeeBeanRepo := repository.NewCoffeeBeanRepository(queries)

	authService := services.NewAuthService(userRepo, jwtManager)
	authHandler := handlers.NewAuthHandler(authService)
	coffeeBeansService := services.NewCoffeeBeansService(coffeeBeanRepo)
	coffeeBeansHandler := handlers.NewCoffeeBeansHandler(coffeeBeansService)

	r := router.New(router.Deps{
		AuthHandler:        authHandler,
		CoffeeBeansHandler: coffeeBeansHandler,
		JWTManager:         jwtManager,
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

	go func() {
		slog.Info("server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}
	slog.Info("server stopped")
}
