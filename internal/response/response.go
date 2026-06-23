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

func DecodeJSON(
	r *http.Request,
	dst interface{},
) error {

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(dst)
}
