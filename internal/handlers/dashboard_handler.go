package handlers

import (
	"net/http"

	"github.com/Hettank/habit-tracker/internal/middleware"
	"github.com/Hettank/habit-tracker/internal/response"
	"github.com/Hettank/habit-tracker/internal/services"
)

type DashboardHandler struct {
	dashboardService *services.DashboardService
}

func NewDashboardHandler(
	dashboardService *services.DashboardService,
) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

func (h *DashboardHandler) GetDashboard(
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
	dashboard, err := h.dashboardService.GetDashboard(
		r.Context(),
		claims.UserID,
	)

	if err != nil {
		response.InternalServerError(
			w,
			"failed to fetch dashboard data",
		)
		return
	}

	response.Success(
		w,
		http.StatusOK,
		"dashboard fetched successfully",
		dashboard,
	)
}
