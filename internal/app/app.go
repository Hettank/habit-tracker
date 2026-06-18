package app

import (
	"github.com/Hettank/habit-tracker/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Config *config.Config
	DB     *pgxpool.Pool
}

func New(cfg *config.Config, db *pgxpool.Pool) *App {

	return &App{
		Config: cfg,
		DB:     db,
	}
}
