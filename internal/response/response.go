package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func Success(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, statusCode int, message string, errors interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

func Created(w http.ResponseWriter, message string, data interface{}) {
	Success(
		w,
		http.StatusCreated,
		message,
		data,
	)
}

func BadRequest(
	w http.ResponseWriter,
	message string,
	errors interface{},
) {
	Error(
		w,
		http.StatusBadRequest,
		message,
		errors,
	)
}

func Unauthorized(
	w http.ResponseWriter,
	message string,
) {
	Error(
		w,
		http.StatusUnauthorized,
		message,
		nil,
	)
}

func Forbidden(
	w http.ResponseWriter,
	message string,
) {
	Error(
		w,
		http.StatusForbidden,
		message,
		nil,
	)
}

func NotFound(
	w http.ResponseWriter,
	message string,
) {
	Error(
		w,
		http.StatusNotFound,
		message,
		nil,
	)
}

func InternalServerError(
	w http.ResponseWriter,
	message string,
) {
	Error(
		w,
		http.StatusInternalServerError,
		message,
		nil,
	)
}

func Conflict(
	w http.ResponseWriter,
	message string,
) {
	Error(
		w,
		http.StatusConflict,
		message,
		nil,
	)
}

func MethodNotAllowed(
	w http.ResponseWriter,
	message string,
) {
	Error(
		w,
		http.StatusMethodNotAllowed,
		message,
		nil,
	)
}

func ValidationError(
	w http.ResponseWriter,
	errors map[string]string,
) {
	Error(
		w,
		http.StatusBadRequest,
		"Validation failed",
		errors,
	)
}
