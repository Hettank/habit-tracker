package services

import (
	"context"

	"github.com/Hettank/habit-tracker/internal/dto"
	apperrors "github.com/Hettank/habit-tracker/internal/errors"
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

func (s *HabitService) GetAll(
	ctx context.Context,
	userID int64,
) ([]models.Habit, error) {
	return s.habitRepo.GetAllByUserID(ctx, userID)
}

func (s *HabitService) GetByID(
	ctx context.Context,
	id int64,
	userID int64,
) (*models.Habit, error) {

	habit, err := s.habitRepo.GetByIDAndUserID(
		ctx,
		id,
		userID,
	)

	if err != nil {
		return nil, err
	}

	if habit == nil {
		return nil, apperrors.ErrNotFound
	}

	return habit, nil
}

func (s *HabitService) Update(
	ctx context.Context,
	id int64,
	userID int64,
	req dto.UpdateHabitRequest,
) (*models.Habit, error) {

	habit, err := s.habitRepo.GetByIDAndUserID(
		ctx,
		id,
		userID,
	)

	if err != nil {
		return nil, err
	}

	if habit == nil {
		return nil, apperrors.ErrNotFound
	}

	habit.Title = req.Title
	habit.Description = req.Description
	habit.Color = req.Color
	habit.Icon = req.Icon

	if err := s.habitRepo.Update(ctx, habit); err != nil {
		return nil, err
	}

	return habit, nil
}

func (s *HabitService) Delete(
	ctx context.Context,
	id int64,
	userID int64,
) error {

	habit, err := s.habitRepo.GetByIDAndUserID(
		ctx,
		id,
		userID,
	)

	if err != nil {
		return err
	}

	if habit == nil {
		return apperrors.ErrNotFound
	}

	return s.habitRepo.Delete(ctx, habit.ID)
}

func (s *HabitService) CheckIn(
	ctx context.Context,
	id int64,
	userID int64,
) error {

	habit, err := s.habitRepo.GetByIDAndUserID(
		ctx,
		id,
		userID,
	)

	if err != nil {
		return err
	}

	if habit == nil {
		return apperrors.ErrNotFound
	}

	exists, err := s.habitRepo.CheckInExistsToday(ctx, habit.ID)
	if err != nil {
		return err
	}

	if exists {
		return apperrors.ErrAlreadyCheckedIn
	}

	return s.habitRepo.CreateCheckIn(ctx, habit.ID)
}

func (s *HabitService) GetCheckedInToday(
	ctx context.Context,
	userID int64,
) ([]models.Habit, error) {
	return s.habitRepo.GetCheckedInHabitsToday(ctx, userID)
}
