package routes

import (
	"net/http"

	"github.com/Hettank/habit-tracker/internal/handlers"
	"github.com/Hettank/habit-tracker/internal/middleware"
	"github.com/Hettank/habit-tracker/internal/utils"
)

func SetupRoutes(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	habitHandler *handlers.HabitHandler,
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
		"POST /api/v1/habits",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(habitHandler.CreateHabit),
		),
	)

	return mux
}
