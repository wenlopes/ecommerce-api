package api

import (
	"encoding/json"
	"net/http"
)

type ErrorBody struct {
	Error string `json:"error"`
}

// OKResponse sends a JSON response with the given status code and data.
func OKResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	_ = json.NewEncoder(w).Encode(data)
}

// ErrorResponse sends a JSON error response with the given status code and message.
func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(ErrorBody{
		Error: message,
	})
}
