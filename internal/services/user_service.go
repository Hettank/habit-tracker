package services

import (
	"context"

	"github.com/Hettank/habit-tracker/internal/models"
	"github.com/Hettank/habit-tracker/internal/repositories"
)

// UserService contains business logic for user profile operations.
type UserService struct {
	userRepo *repositories.UserRepository
}

// NewUserservice creates a new UserService with the given repository dependency.
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
