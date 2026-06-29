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
