package routes

import (
	"net/http"

	"github.com/Hettank/habit-tracker/internal/handlers"
	"github.com/Hettank/habit-tracker/internal/middleware"
	"github.com/Hettank/habit-tracker/internal/utils"
)

// SetupRoutes registers all API routes and applies authentication middleware.
func SetupRoutes(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	habitHandler *handlers.HabitHandler,
	dashboardHandler *handlers.DashboardHandler,
	jwtManager *utils.JWTManager,
) *http.ServeMux {

	mux := http.NewServeMux()

	// ================== Auth Routes ==================

	mux.HandleFunc(
		"GET /health",
		handlers.Health,
	)

	mux.HandleFunc(
		"POST /api/v1/auth/register",
		authHandler.Register,
	)

	mux.HandleFunc(
		"POST /api/v1/auth/login",
		authHandler.Login,
	)

	mux.Handle(
		"POST /api/v1/auth/logout",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(authHandler.Logout),
		),
	)

	mux.Handle(
		"POST /api/v1/auth/logout-all",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(authHandler.LogoutAll),
		),
	)

	mux.HandleFunc(
		"POST /api/v1/auth/refresh",
		authHandler.Refresh,
	)

	mux.Handle(
		"GET /api/v1/me",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(userHandler.Me),
		),
	)

	// ================== Habit Routes ==================

	mux.Handle(
		"GET /api/v1/habits/today",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(habitHandler.GetCheckedInHabitsToday),
		),
	)

	mux.Handle(
		"GET /api/v1/habits/{id}/streak",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(habitHandler.GetHabitStreak),
		),
	)

	mux.Handle(
		"GET /api/v1/habits/{id}/history",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(habitHandler.GetHabitHistory),
		),
	)

	mux.Handle(
		"POST /api/v1/habits/{id}/check-in",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(habitHandler.CheckInHabit),
		),
	)

	mux.Handle(
		"GET /api/v1/habits/{id}",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(habitHandler.GetHabitByID),
		),
	)

	mux.Handle(
		"PUT /api/v1/habits/{id}",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(habitHandler.UpdateHabit),
		),
	)

	mux.Handle(
		"DELETE /api/v1/habits/{id}",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(habitHandler.DeleteHabit),
		),
	)

	mux.Handle(
		"POST /api/v1/habits",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(habitHandler.CreateHabit),
		),
	)

	mux.Handle(
		"GET /api/v1/habits",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(habitHandler.GetHabits),
		),
	)

	// ================== Dashboard Routes ==================
	mux.Handle(
		"GET /api/v1/dashboard",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(dashboardHandler.GetDashboard),
		),
	)

	return mux
}
