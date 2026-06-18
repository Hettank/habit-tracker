package routes

import (
	"net/http"

	"github.com/Hettank/habit-tracker/internal/handlers"
)

func SetupRoutes() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /health",
		handlers.Health,
	)

	return mux
}
