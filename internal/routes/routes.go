package routes

import (
	"net/http"

	"github.com/Hettank/habit-tracker/internal/handlers"
)

func SetupRoutes(
	authHandler *handlers.AuthHandler,
) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /health",
		handlers.Health,
	)

	mux.HandleFunc(
		"POST /api/v1/auth/register",
		authHandler.Register,
	)

	return mux
}
