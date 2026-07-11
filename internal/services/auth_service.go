package services

import (
	"context"
	"errors"
	"time"

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

func (s *AuthService) createSession(
	ctx context.Context,
	user *models.User,
) (*dto.AuthResult, error) {

	accessToken, err := s.jwtManager.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshModel := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: utils.HashToken(refreshToken),
		ExpiresAt: time.Now().Add(
			s.jwtManager.RefreshTokenTTL(),
		),
	}

	if err := s.refreshRepo.Create(
		ctx,
		refreshModel,
	); err != nil {
		return nil, err
	}

	return &dto.AuthResult{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		RefreshExpiry: refreshModel.ExpiresAt,
	}, nil
}

func (s *AuthService) Register(
	ctx context.Context,
	req dto.RegisterRequest,
) (*models.User, error) {
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

func (s *AuthService) Login(
	ctx context.Context,
	req dto.LoginRequest,
) (*dto.AuthResult, error) {

	user, err := s.userRepo.GetByEmail(
		ctx,
		req.Email,
	)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	if err := utils.CheckPassword(
		user.PasswordHash,
		req.Password,
	); err != nil {

		return nil, apperrors.ErrInvalidCredentials
	}

	return s.createSession(ctx, user)
}

func (s *AuthService) Refresh(
	ctx context.Context,
	refreshToken string,
) (*dto.AuthResult, error) {

	// Step 1: Hash incoming refresh token
	hashedToken := utils.HashToken(refreshToken)

	// Step 2: Find token in database
	tokenModel, err := s.refreshRepo.GetByTokenHash(
		ctx,
		hashedToken,
	)
	if err != nil {
		return nil, err
	}

	if tokenModel == nil {
		return nil, apperrors.ErrUnauthorized
	}

	// Step 3: Check if token has expired
	if time.Now().After(tokenModel.ExpiresAt) {

		// Cleanup expired token
		_ = s.refreshRepo.DeleteById(
			ctx,
			tokenModel.ID,
		)

		return nil, apperrors.ErrUnauthorized
	}

	// Optional (when you start using RevokedAt)
	// if tokenModel.RevokedAt != nil {
	//     return nil, apperrors.ErrUnauthorized
	// }

	// Step 4: Load user
	user, err := s.userRepo.GetByID(
		ctx,
		tokenModel.UserID,
	)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, apperrors.ErrUnauthorized
	}

	// Step 5: Delete old refresh token
	if err := s.refreshRepo.DeleteById(
		ctx,
		tokenModel.ID,
	); err != nil {
		return nil, err
	}

	return s.createSession(ctx, user)
}

func (s *AuthService) Logout(
	ctx context.Context,
	refreshToken string,
) error {
	hashedToken := utils.HashToken(refreshToken)

	err := s.refreshRepo.RevokeByTokenHash(ctx, hashedToken)

	if errors.Is(err, apperrors.ErrRefreshTokenNotFound) {
		return apperrors.ErrUnauthorized
	}

	return err
}

func (s *AuthService) LogoutAll(
	ctx context.Context,
	userID int64,
) error {
	return s.refreshRepo.RevokeAllByUserId(
		ctx,
		userID,
	)
}
