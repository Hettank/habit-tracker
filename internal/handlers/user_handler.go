package handlers

import (
	"net/http"

	"github.com/Hettank/habit-tracker/internal/dto"
	"github.com/Hettank/habit-tracker/internal/middleware"
	"github.com/Hettank/habit-tracker/internal/response"
	"github.com/Hettank/habit-tracker/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(
	userService *services.UserService,
) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Me(
	w http.ResponseWriter,
	r *http.Request,
) {
	claims, ok := middleware.GetClaims(r.Context())

	if !ok {
		response.Unauthorized(
			w,
			"unauthorized",
		)
		return
	}

	user, err := h.userService.Me(
		r.Context(),
		claims.UserID,
	)

	if err != nil {
		response.InternalServerError(
			w,
			"failed to fetch user",
		)
		return
	}

	response.Success(
		w,
		http.StatusOK,
		"user fetched successfully",
		dto.MeResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	)
}
