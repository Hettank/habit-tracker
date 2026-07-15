package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Hettank/habit-tracker/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BuildDSN constructs a PostgreSQL connection string from the application configuration.
func BuildDSN(cfg *config.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)
}

// New creates a PostgreSQL connection pool and verifies connectivity with a ping.
func New(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := BuildDSN(cfg)
	log.Println("DSN:", dsn)
	pool, err := pgxpool.New(context.Background(), dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Improved Ping with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close() // Clean up the pool if ping fails
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

// Close gracefully shuts down the PostgreSQL connection pool.
func Close(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
		log.Println("Database connection pool closed")
	}
}
