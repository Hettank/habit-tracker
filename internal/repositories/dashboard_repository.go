package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DashboardRepository provides data access methods for dashboard statistics.
type DashboardRepository struct {
	db *pgxpool.Pool
}

// NewDashboardRepository creates a new DashboardRepository with the given connection pool.
func NewDashboardRepository(db *pgxpool.Pool) *DashboardRepository {
	return &DashboardRepository{
		db: db,
	}
}

// DashboardStats holds raw counts fetched from the database for dashboard calculations.
type DashboardStats struct {
	TotalHabits    int64
	CompletedToday int64
}

func (r *DashboardRepository) GetUserStats(
	ctx context.Context,
	userID int64,
) (*DashboardStats, error) {
	var totalHabits int64
	queryTotal := `
		SELECT COUNT(id)
		FROM habits
		WHERE user_id = $1 AND deleted_at IS NULL
	`
	err := r.db.QueryRow(ctx, queryTotal, userID).Scan(&totalHabits)
	if err != nil {
		return nil, err
	}

	var completedToday int64
	queryCompleted := `
		SELECT COUNT(DISTINCT hl.habit_id)
		FROM habit_logs hl
		JOIN habits h ON hl.habit_id = h.id
		WHERE h.user_id = $1
			AND h.deleted_at IS NULL
			AND DATE(hl.completed_at) = CURRENT_DATE
	`
	err = r.db.QueryRow(ctx, queryCompleted, userID).Scan(&completedToday)
	if err != nil {
		return nil, err
	}

	return &DashboardStats{
		TotalHabits:    totalHabits,
		CompletedToday: completedToday,
	}, nil
}
