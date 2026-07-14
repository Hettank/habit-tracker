package repositories

import (
	"context"
	"errors"

	"github.com/Hettank/habit-tracker/internal/models"
	"github.com/jackc/pgx/v5"
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

func (r *HabitRepository) GetAllByUserID(
	ctx context.Context,
	userID int64,
) ([]models.Habit, error) {

	query := `
		SELECT
			id,
			user_id,
			title,
			description,
			color,
			icon,
			created_at,
			updated_at
		FROM habits
		WHERE user_id = $1
			AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habits []models.Habit

	for rows.Next() {
		var habit models.Habit

		if err := rows.Scan(
			&habit.ID,
			&habit.UserID,
			&habit.Title,
			&habit.Description,
			&habit.Color,
			&habit.Icon,
			&habit.CreatedAt,
			&habit.UpdatedAt,
		); err != nil {
			return nil, err
		}

		habits = append(habits, habit)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return habits, nil
}

func (r *HabitRepository) GetByID(
	ctx context.Context,
	id int64,
) (*models.Habit, error) {

	query := `
		SELECT
			id,
			user_id,
			title,
			description,
			color,
			icon,
			created_at,
			updated_at
		FROM habits
		WHERE id = $1
			AND deleted_at IS NULL
	`

	var habit models.Habit

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&habit.ID,
		&habit.UserID,
		&habit.Title,
		&habit.Description,
		&habit.Color,
		&habit.Icon,
		&habit.CreatedAt,
		&habit.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &habit, nil
}

func (r *HabitRepository) GetByIDAndUserID(
	ctx context.Context,
	id int64,
	userID int64,
) (*models.Habit, error) {

	query := `
		SELECT
			id,
			user_id,
			title,
			description,
			color,
			icon,
			created_at,
			updated_at
		FROM habits
		WHERE id = $1
			AND user_id = $2
			AND deleted_at IS NULL
	`

	var habit models.Habit

	err := r.db.QueryRow(
		ctx,
		query,
		id,
		userID,
	).Scan(
		&habit.ID,
		&habit.UserID,
		&habit.Title,
		&habit.Description,
		&habit.Color,
		&habit.Icon,
		&habit.CreatedAt,
		&habit.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &habit, nil
}

func (r *HabitRepository) Update(
	ctx context.Context,
	habit *models.Habit,
) error {

	query := `
		UPDATE habits
		SET
			title = $1,
			description = $2,
			color = $3,
			icon = $4,
			updated_at = NOW()
		WHERE id = $5
			AND deleted_at IS NULL
		RETURNING
			updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		habit.Title,
		habit.Description,
		habit.Color,
		habit.Icon,
		habit.ID,
	).Scan(
		&habit.UpdatedAt,
	)
}

func (r *HabitRepository) Delete(
	ctx context.Context,
	id int64,
) error {

	query := `
		UPDATE habits
		SET deleted_at = NOW()
		WHERE id = $1
			AND deleted_at IS NULL
	`

	_, err := r.db.Exec(
		ctx,
		query,
		id,
	)

	return err
}

func (r *HabitRepository) CheckInExistsToday(
	ctx context.Context,
	habitID int64,
) (bool, error) {

	query := `
		SELECT EXISTS(
			SELECT 1 FROM habit_logs
			WHERE habit_id = $1 AND DATE(completed_at) = CURRENT_DATE
		)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, habitID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *HabitRepository) CreateCheckIn(
	ctx context.Context,
	habitID int64,
) error {
	query := `
		INSERT INTO habit_logs (habit_id, completed_at)
		VALUES ($1, NOW())
	`
	_, err := r.db.Exec(ctx, query, habitID)
	return err
}

func (r *HabitRepository) GetCheckedInHabitsToday(
	ctx context.Context,
	userID int64,
) ([]models.Habit, error) {

	query := `
		SELECT
			h.id,
			h.user_id,
			h.title,
			h.description,
			h.color,
			h.icon,
			h.created_at,
			h.updated_at
		FROM habits h
		INNER JOIN habit_logs hl ON h.id = hl.habit_id
		WHERE h.user_id = $1
			AND h.deleted_at IS NULL
			AND DATE(hl.completed_at) = CURRENT_DATE
		ORDER BY hl.completed_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habits []models.Habit

	for rows.Next() {
		var habit models.Habit

		if err := rows.Scan(
			&habit.ID,
			&habit.UserID,
			&habit.Title,
			&habit.Description,
			&habit.Color,
			&habit.Icon,
			&habit.CreatedAt,
			&habit.UpdatedAt,
		); err != nil {
			return nil, err
		}

		habits = append(habits, habit)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return habits, nil
}

func (r *HabitRepository) GetHabitHistory(
	ctx context.Context,
	habitID int64,
) ([]models.HabitLog, error) {
	query := `
		SELECT
			id,
			habit_id,
			completed_at,
			created_at
		FROM habit_logs
		WHERE habit_id = $1
		ORDER BY completed_at DESC
	`

	rows, err := r.db.Query(ctx, query, habitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.HabitLog

	for rows.Next() {
		var log models.HabitLog
		if err := rows.Scan(
			&log.ID,
			&log.HabitID,
			&log.CompletedAt,
			&log.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

