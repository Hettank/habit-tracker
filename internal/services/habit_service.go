package services

import (
	"context"

	"github.com/Hettank/habit-tracker/internal/dto"
	"github.com/Hettank/habit-tracker/internal/models"
	"github.com/Hettank/habit-tracker/internal/repositories"
)

type HabitService struct {
	habitRepo *repositories.HabitRepository
}

func NewHabitService(
	habitRepo *repositories.HabitRepository,
) *HabitService {
	return &HabitService{
		habitRepo: habitRepo,
	}
}

func (s *HabitService) Create(
	ctx context.Context,
	userID int64,
	req dto.CreateHabitRequest,
) (*models.Habit, error) {

	habit := &models.Habit{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
	}

	if err := s.habitRepo.Create(ctx, habit); err != nil {
		return nil, err
	}

	return habit, nil
}
