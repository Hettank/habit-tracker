package services

import (
	"context"

	"github.com/Hettank/habit-tracker/internal/dto"
	apperrors "github.com/Hettank/habit-tracker/internal/errors"
	"github.com/Hettank/habit-tracker/internal/models"
	"github.com/Hettank/habit-tracker/internal/repositories"
	"github.com/Hettank/habit-tracker/internal/utils"
)

// whenever the service needs user data, it asks the repository.
type AuthService struct {
	userRepo    *repositories.UserRepository
	refreshRepo *repositories.RefreshTokenRepository
	jwtManager  *utils.JWTManager
}

func NewAuthService(
	userRepo *repositories.UserRepository,
	refreshRepo *repositories.RefreshTokenRepository,
	jwtManager *utils.JWTManager,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
		jwtManager:  jwtManager,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*models.User, error) {
	// Step 1: Check if email already exists.
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)

	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, apperrors.ErrUserAlreadyExists
	}

	// Step 2: Hash password.
	hashedPassword, err := utils.HashPassword(
		req.Password,
	)

	if err != nil {
		return nil, err
	}

	// Step 3: Create domain model
	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	// Step 4: save user
	err = s.userRepo.Create(ctx, user)

	if err != nil {
		return nil, err
	}

	// return the created user.
	return user, nil
}
