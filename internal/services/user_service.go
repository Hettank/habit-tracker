package services

import (
	"context"

	"github.com/Hettank/habit-tracker/internal/models"
	"github.com/Hettank/habit-tracker/internal/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserservice(
	userRepo *repositories.UserRepository,
) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Me(
	ctx context.Context,
	userID int64,
) (*models.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}
