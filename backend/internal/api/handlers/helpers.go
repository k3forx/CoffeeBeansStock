package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/api"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api/middleware"
)

// requireUserID extracts the authenticated user ID from the request context.
// It writes a 401 response and returns false if the user is not authenticated.
func requireUserID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		api.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "認証が必要です")
	}
	return userID, ok
}

// decodeJSON decodes the request body into dst.
// It writes a 400 response and returns false on failure.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		api.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "リクエストの形式が不正です")
		return false
	}
	return true
}

// parseUUIDParam parses a UUID from a URL path parameter.
// It writes a 400 response and returns false on failure.
func parseUUIDParam(w http.ResponseWriter, r *http.Request, param string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, param))
	if err != nil {
		api.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "不正なIDです")
		return uuid.Nil, false
	}
	return id, true
}

// toStringPtr converts a pointer to a string-based type to a *string.
func toStringPtr[T ~string](v *T) *string {
	if v == nil {
		return nil
	}
	s := string(*v)
	return &s
}

// parseQueryInt parses an integer query parameter with a default value.
// The result is clamped to [min, max].
func parseQueryInt(r *http.Request, key string, defaultVal, min, max int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}
