package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Response is the standard API response envelope.
type Response struct {
	Data    any        `json:"data,omitempty"`
	Error   *ErrorBody `json:"error,omitempty"`
	Success bool       `json:"success"`
}

// ErrorBody contains error details in the API response.
type ErrorBody struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

// FieldError represents a validation error for a specific field.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// WriteJSON writes a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}

// WriteSuccess writes a successful JSON response.
func WriteSuccess(w http.ResponseWriter, status int, data any) {
	WriteJSON(w, status, Response{Success: true, Data: data, Error: nil})
}

// WriteError writes an error JSON response.
// For 5xx errors, it also logs the error details server-side.
func WriteError(w http.ResponseWriter, status int, code, message string) {
	if status >= 500 {
		slog.Error("server error",
			"status", status,
			"code", code,
			"message", message,
		)
	}
	WriteJSON(w, status, Response{
		Success: false,
		Data:    nil,
		Error:   &ErrorBody{Code: code, Message: message, Details: nil},
	})
}

// WriteValidationError writes a validation error response with field-level details.
func WriteValidationError(w http.ResponseWriter, details []FieldError) {
	WriteJSON(w, http.StatusBadRequest, Response{
		Success: false,
		Data:    nil,
		Error: &ErrorBody{
			Code:    "VALIDATION_ERROR",
			Message: "入力値が不正です",
			Details: details,
		},
	})
}
