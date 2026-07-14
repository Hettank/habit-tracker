package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hettank/habit-tracker/internal/app"
	"github.com/Hettank/habit-tracker/internal/config"
	"github.com/Hettank/habit-tracker/internal/db"
	"github.com/Hettank/habit-tracker/internal/handlers"
	"github.com/Hettank/habit-tracker/internal/repositories"
	"github.com/Hettank/habit-tracker/internal/routes"
	"github.com/Hettank/habit-tracker/internal/services"
	"github.com/Hettank/habit-tracker/internal/utils"
)

func main() {
	// 1. Load Configuration
	cfg := config.Load()

	// 2. Initialize Dependencies (Database)
	dbPool, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close(dbPool)

	log.Println("Database Connected successfully.")

	// Dependency Injection
	userRepo := repositories.NewUserRepository(dbPool)
	refreshRepo := repositories.NewRefreshTokenRepository(dbPool)
	habitRepo := repositories.NewHabitRepository(dbPool)

	jwtManager := utils.NewJWTManager(
		cfg.JWTSecret,
		cfg.AccessTokenTTL,
		cfg.RefreshTokenTTL,
	)

	authService := services.NewAuthService(
		userRepo,
		refreshRepo,
		jwtManager,
	)

	userService := services.NewUserservice(userRepo)
	habitService := services.NewHabitService(habitRepo)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	habitHandler := handlers.NewHabitHandler(habitService)

	mux := routes.SetupRoutes(
		authHandler,
		userHandler,
		habitHandler,
		jwtManager,
	)

	// 3. Create App (Dependency Injection container)
	application := app.New(cfg, dbPool)
	_ = application

	// Create HTTP Server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: mux,
	}

	// Create a Goroutine
	go func() {
		log.Printf("Server Starting on http://localhost:%s", cfg.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server Failed: %v", err)
		}
	}()

	// Create a channel to receive OS signals (Ctrl+C, kill, etc.)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block (wait) here until we receive a shutdown signal
	<-quit
	log.Println("Shutting down server...")

	// Give ongoing requests time to finish (Graceful Shutdown)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server exited gracefully")
	}
}
