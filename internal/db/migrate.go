package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// RunMigrations connects to the database and applies all pending schema migrations.
// It returns nil if migrations succeed or if no pending migrations are found.
func RunMigrations(dsn string) error {
	log.Println("Running database migrations...")

	// Create a database/sql connection using pgx driver.
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database for migrations: %w", err)
	}
	defer db.Close()

	// Verify the database connection with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database for migrations: %w", err)
	}

	// Create the postgres migration driver.
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Create the migrate instance.
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate instance: %w", err)
	}

	// Clean up resources.
	defer func() {
		sourceErr, dbErr := m.Close()

		if sourceErr != nil {
			log.Printf("Error closing migration source: %v", sourceErr)
		}

		if dbErr != nil {
			log.Printf("Error closing migration database: %v", dbErr)
		}
	}()

	// Run all pending migrations.
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No pending migrations found.")
			return nil
		}

		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully.")

	return nil
}
