package handlers

import (
	"net/http"

	"github.com/Hettank/habit-tracker/internal/response"
)

func Health(w http.ResponseWriter, r *http.Request) {

	response.Success(
		w,
		http.StatusOK,
		"Server Healthy",
		map[string]string{
			"status": "ok",
		},
	)
}
