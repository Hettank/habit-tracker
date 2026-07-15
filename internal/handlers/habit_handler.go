package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Hettank/habit-tracker/internal/dto"
	apperrors "github.com/Hettank/habit-tracker/internal/errors"
	"github.com/Hettank/habit-tracker/internal/middleware"
	"github.com/Hettank/habit-tracker/internal/models"
	"github.com/Hettank/habit-tracker/internal/response"
	"github.com/Hettank/habit-tracker/internal/services"
	"github.com/Hettank/habit-tracker/internal/validator"
)

// HabitHandler handles HTTP requests for habit-related endpoints.
type HabitHandler struct {
	habitService *services.HabitService
}

// NewHabitHandler creates a new HabitHandler with the given service dependency.
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

	response.Created(
		w,
		"habit created successfully",
		toHabitResponse(habit),
	)
}

func (h *HabitHandler) GetHabits(
	w http.ResponseWriter,
	r *http.Request,
) {
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
	habits, err := h.habitService.GetAll(
		r.Context(),
		claims.UserID,
	)

	if err != nil {
		response.InternalServerError(
			w,
			"failed to fetch habits",
		)
		return
	}

	res := make([]dto.HabitResponse, len(habits))
	for i, habit := range habits {
		res[i] = toHabitResponse(&habit)
	}

	response.Success(
		w,
		http.StatusOK,
		"habits fetched successfully",
		res,
	)
}

func (h *HabitHandler) GetHabitByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	// Parse habit ID from URL
	id, err := parseIDParam(r)
	if err != nil {
		response.BadRequest(
			w,
			"invalid habit id",
			nil,
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
	habit, err := h.habitService.GetByID(
		r.Context(),
		id,
		claims.UserID,
	)

	if err != nil {

		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			response.NotFound(
				w,
				"habit not found",
			)
		default:
			response.InternalServerError(
				w,
				"failed to fetch habit",
			)
		}

		return
	}

	response.Success(
		w,
		http.StatusOK,
		"habit fetched successfully",
		toHabitResponse(habit),
	)
}

func (h *HabitHandler) UpdateHabit(
	w http.ResponseWriter,
	r *http.Request,
) {
	// Parse habit ID from URL
	id, err := parseIDParam(r)
	if err != nil {
		response.BadRequest(
			w,
			"invalid habit id",
			nil,
		)
		return
	}

	var req dto.UpdateHabitRequest

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
	habit, err := h.habitService.Update(
		r.Context(),
		id,
		claims.UserID,
		req,
	)

	if err != nil {

		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			response.NotFound(
				w,
				"habit not found",
			)
		default:
			response.InternalServerError(
				w,
				"failed to update habit",
			)
		}

		return
	}

	response.Success(
		w,
		http.StatusOK,
		"habit updated successfully",
		toHabitResponse(habit),
	)
}

func (h *HabitHandler) DeleteHabit(
	w http.ResponseWriter,
	r *http.Request,
) {
	// Parse habit ID from URL
	id, err := parseIDParam(r)
	if err != nil {
		response.BadRequest(
			w,
			"invalid habit id",
			nil,
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
	err = h.habitService.Delete(
		r.Context(),
		id,
		claims.UserID,
	)

	if err != nil {

		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			response.NotFound(
				w,
				"habit not found",
			)
		default:
			response.InternalServerError(
				w,
				"failed to delete habit",
			)
		}

		return
	}

	response.Success(
		w,
		http.StatusOK,
		"habit deleted successfully",
		nil,
	)
}

func (h *HabitHandler) CheckInHabit(
	w http.ResponseWriter,
	r *http.Request,
) {
	// Parse habit ID from URL
	id, err := parseIDParam(r)
	if err != nil {
		response.BadRequest(
			w,
			"invalid habit id",
			nil,
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
	err = h.habitService.CheckIn(
		r.Context(),
		id,
		claims.UserID,
	)

	if err != nil {

		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			response.NotFound(
				w,
				"habit not found",
			)
		case errors.Is(err, apperrors.ErrAlreadyCheckedIn):
			response.Conflict(
				w,
				err.Error(),
			)
		default:
			response.InternalServerError(
				w,
				"failed to check in habit",
			)
		}

		return
	}

	response.Created(
		w,
		"habit checked in successfully",
		nil,
	)
}

func (h *HabitHandler) GetCheckedInHabitsToday(
	w http.ResponseWriter,
	r *http.Request,
) {
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
	habits, err := h.habitService.GetCheckedInToday(
		r.Context(),
		claims.UserID,
	)

	if err != nil {
		response.InternalServerError(
			w,
			"failed to fetch checked in habits",
		)
		return
	}

	res := make([]dto.HabitResponse, len(habits))
	for i, habit := range habits {
		res[i] = toHabitResponse(&habit)
	}

	response.Success(
		w,
		http.StatusOK,
		"checked in habits fetched successfully",
		res,
	)
}

func (h *HabitHandler) GetHabitHistory(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := parseIDParam(r)
	if err != nil {
		response.BadRequest(
			w,
			"invalid habit id",
			nil,
		)
		return
	}

	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		response.Unauthorized(
			w,
			"unauthorized",
		)
		return
	}

	logs, err := h.habitService.GetHistory(
		r.Context(),
		id,
		claims.UserID,
	)

	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			response.NotFound(
				w,
				"habit not found",
			)
		default:
			response.InternalServerError(
				w,
				"failed to fetch habit history",
			)
		}
		return
	}

	res := make([]dto.HabitLogResponse, len(logs))
	for i, log := range logs {
		res[i] = dto.HabitLogResponse{
			ID:          log.ID,
			HabitID:     log.HabitID,
			CompletedAt: log.CompletedAt,
			CreatedAt:   log.CreatedAt,
		}
	}

	response.Success(
		w,
		http.StatusOK,
		"habit history fetched successfully",
		res,
	)
}

func (h *HabitHandler) GetHabitStreak(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := parseIDParam(r)
	if err != nil {
		response.BadRequest(
			w,
			"invalid habit id",
			nil,
		)
		return
	}

	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		response.Unauthorized(
			w,
			"unauthorized",
		)
		return
	}

	streak, err := h.habitService.GetStreak(
		r.Context(),
		id,
		claims.UserID,
	)

	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			response.NotFound(
				w,
				"habit not found",
			)
		default:
			response.InternalServerError(
				w,
				"failed to fetch habit streak",
			)
		}
		return
	}

	res := dto.HabitStreakResponse{
		CurrentStreak: streak,
	}

	response.Success(
		w,
		http.StatusOK,
		"habit streak fetched successfully",
		res,
	)
}

// parseIDParam extracts the {id} path value from the request URL.
func parseIDParam(r *http.Request) (int64, error) {
	return strconv.ParseInt(
		r.PathValue("id"),
		10,
		64,
	)
}

// toHabitResponse converts a domain model to a DTO response.
func toHabitResponse(habit *models.Habit) dto.HabitResponse {
	return dto.HabitResponse{
		ID:          habit.ID,
		Title:       habit.Title,
		Description: habit.Description,
		Color:       habit.Color,
		Icon:        habit.Icon,
		CreatedAt:   habit.CreatedAt,
		UpdatedAt:   habit.UpdatedAt,
	}
}
