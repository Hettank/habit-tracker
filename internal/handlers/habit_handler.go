package handlers

import (
	"net/http"

	"github.com/Hettank/habit-tracker/internal/dto"
	"github.com/Hettank/habit-tracker/internal/middleware"
	"github.com/Hettank/habit-tracker/internal/response"
	"github.com/Hettank/habit-tracker/internal/services"
	"github.com/Hettank/habit-tracker/internal/validator"
)

type HabitHandler struct {
	habitService *services.HabitService
}

func NewHabitHandler(
	habitService *services.HabitService,
) *HabitHandler {
	return &HabitHandler{
		habitService: habitService,
	}
}

func (h *HabitHandler) CreateHabit(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req dto.CreateHabitRequest

	// Decode JSON
	if err := response.DecodeJSON(r, &req); err != nil {
		response.BadRequest(
			w,
			"invalid request body",
			nil,
		)
		return
	}

	// Validate
	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(
			w,
			validator.FormatValidationErrors(err),
		)
		return
	}

	// Get authenticated user
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		response.Unauthorized(
			w,
			"unauthorized",
		)
		return
	}

	// Call service
	habit, err := h.habitService.Create(
		r.Context(),
		claims.UserID,
		req,
	)

	if err != nil {
		response.InternalServerError(
			w,
			"failed to create habit",
		)
		return
	}

	res := dto.HabitResponse{
		ID:          habit.ID,
		Title:       habit.Title,
		Description: habit.Description,
		Color:       habit.Color,
		Icon:        habit.Icon,
		CreatedAt:   habit.CreatedAt,
		UpdatedAt:   habit.UpdatedAt,
	}

	response.Created(
		w,
		"habit created successfully",
		res,
	)
}
