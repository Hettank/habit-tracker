package services

import (
	"context"

	"github.com/Hettank/habit-tracker/internal/dto"
	"github.com/Hettank/habit-tracker/internal/repositories"
)

// DashboardService contains business logic for calculating dashboard statistics.
type DashboardService struct {
	dashboardRepo *repositories.DashboardRepository
}

// NewDashboardService creates a new DashboardService with the given repository dependency.
func NewDashboardService(
	dashboardRepo *repositories.DashboardRepository,
) *DashboardService {
	return &DashboardService{
		dashboardRepo: dashboardRepo,
	}
}

func (s *DashboardService) GetDashboard(
	ctx context.Context,
	userID int64,
) (*dto.DashboardResponse, error) {
	stats, err := s.dashboardRepo.GetUserStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	pendingToday := stats.TotalHabits - stats.CompletedToday
	if pendingToday < 0 {
		pendingToday = 0
	}

	var completionPercentage int64 = 0
	if stats.TotalHabits > 0 {
		completionPercentage = (stats.CompletedToday * 100) / stats.TotalHabits
	}

	return &dto.DashboardResponse{
		TotalHabits:          stats.TotalHabits,
		CompletedToday:       stats.CompletedToday,
		PendingToday:         pendingToday,
		CompletionPercentage: completionPercentage,
	}, nil
}
