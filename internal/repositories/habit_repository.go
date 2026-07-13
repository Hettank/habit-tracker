package repositories

import (
	"context"

	"github.com/Hettank/habit-tracker/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HabitRepository struct {
	db *pgxpool.Pool
}

func NewHabitRepository(db *pgxpool.Pool) *HabitRepository {
	return &HabitRepository{
		db: db,
	}
}

func (r *HabitRepository) Create(
	ctx context.Context,
	habit *models.Habit,
) error {

	query := `
		INSERT INTO habits (
			user_id,
			title,
			description,
			color,
			icon
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		)
		RETURNING
			id,
			created_at,
			updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		habit.UserID,
		habit.Title,
		habit.Description,
		habit.Color,
		habit.Icon,
	).Scan(
		&habit.ID,
		&habit.CreatedAt,
		&habit.UpdatedAt,
	)
}
