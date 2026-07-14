package services

import (
	"context"
	"time"

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

func (s *HabitService) GetHistory(
	ctx context.Context,
	id int64,
	userID int64,
) ([]models.HabitLog, error) {

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

	return s.habitRepo.GetHabitHistory(ctx, habit.ID)
}

func (s *HabitService) GetStreak(
	ctx context.Context,
	id int64,
	userID int64,
) (int64, error) {

	habit, err := s.habitRepo.GetByIDAndUserID(
		ctx,
		id,
		userID,
	)

	if err != nil {
		return 0, err
	}

	if habit == nil {
		return 0, apperrors.ErrNotFound
	}

	logs, err := s.habitRepo.GetHabitHistory(ctx, habit.ID)
	if err != nil {
		return 0, err
	}

	if len(logs) == 0 {
		return 0, nil
	}

	// Helper to extract just the date part (ignoring time)
	truncateToDay := func(t time.Time) time.Time {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	}

	var uniqueDates []time.Time
	for _, log := range logs {
		day := truncateToDay(log.CompletedAt)
		// Since it's ordered by completed_at DESC, the same day will be consecutive.
		if len(uniqueDates) == 0 || !uniqueDates[len(uniqueDates)-1].Equal(day) {
			uniqueDates = append(uniqueDates, day)
		}
	}

	if len(uniqueDates) == 0 {
		return 0, nil
	}

	latest := uniqueDates[0]
	today := truncateToDay(time.Now())
	yesterday := today.AddDate(0, 0, -1)

	// If latest completion is neither today nor yesterday, streak is broken.
	if !latest.Equal(today) && !latest.Equal(yesterday) {
		return 0, nil
	}

	var streak int64 = 1
	for i := 1; i < len(uniqueDates); i++ {
		prev := uniqueDates[i-1]
		curr := uniqueDates[i]

		// The difference must be exactly 1 day.
		if prev.Sub(curr) == 24*time.Hour {
			streak++
		} else {
			break
		}
	}

	return streak, nil
}

