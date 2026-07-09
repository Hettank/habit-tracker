package handlers

import (
	"errors"
	"net/http"

	"github.com/Hettank/habit-tracker/internal/dto"
	apperrors "github.com/Hettank/habit-tracker/internal/errors"
	"github.com/Hettank/habit-tracker/internal/response"
	"github.com/Hettank/habit-tracker/internal/services"
	"github.com/Hettank/habit-tracker/internal/validator"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	// Decode Request body
	if err := response.DecodeJSON(r, &req); err != nil {
		response.BadRequest(
			w,
			"Invalid request body",
			nil,
		)
		return
	}

	// Validate request
	if err := validator.Validate.Struct(req); err != nil {
		response.BadRequest(
			w,
			"validation failed",
			validator.FormatValidationErrors(err),
		)
		return
	}

	// Call service
	user, err := h.authService.Register(
		r.Context(),
		req,
	)

	if err != nil {

		switch {
		case errors.Is(err, apperrors.ErrUserAlreadyExists):

			response.Conflict(
				w,
				err.Error(),
			)
		default:
			response.InternalServerError(
				w,
				"internal server error",
			)
		}

		return
	}

	res := dto.RegisterResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	response.Created(
		w,
		"user registered successfully",
		res,
	)
}

func (h *AuthHandler) Login(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req dto.LoginRequest

	if err := response.DecodeJSON(r, &req); err != nil {
		response.BadRequest(
			w,
			"Invalid request body",
			nil,
		)
		return
	}

	if errs := validator.Validate.Struct(req); errs != nil {
		response.ValidationError(
			w,
			validator.FormatValidationErrors(errs),
		)
		return
	}

	result, err := h.authService.Login(
		r.Context(),
		req,
	)

	if err != nil {

		switch {

		case errors.Is(err, apperrors.ErrInvalidCredentials):
			response.Unauthorized(
				w,
				err.Error(),
			)

		default:
			response.InternalServerError(
				w,
				"internal server error",
			)
		}

		return
	}

	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "refresh_token",
			Value:    result.RefreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			Expires:  result.RefreshExpiry,
		},
	)

	loginResponse := dto.LoginResponse{
		AccessToken: result.AccessToken,
	}

	response.Success(
		w,
		http.StatusOK,
		"Login successful",
		loginResponse,
	)
}

func (h *AuthHandler) Refresh(
	w http.ResponseWriter,
	r *http.Request,
) {
	// Step 1: Read refresh token cookie
	cookie, err := r.Cookie("refresh_token")

	if err != nil {
		response.Unauthorized(
			w,
			"refresh token missing",
		)
		return
	}

	// Step 2: Call Service
	result, err := h.authService.Refresh(
		r.Context(),
		cookie.Value,
	)

	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUnauthorized):
			response.Unauthorized(
				w,
				err.Error(),
			)
		default:
			response.InternalServerError(
				w,
				"Internal Server Error",
			)
		}
	}

	// Step 3: Replace refresh token cookie
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "refresh_token",
			Value:    result.RefreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			Expires:  result.RefreshExpiry,
		},
	)

	// Step 4: Return new access token
	response.Success(
		w,
		http.StatusOK,
		"token refreshed successfully",
		dto.RefreshResponse{
			AccessToken: result.AccessToken,
		},
	)
}
