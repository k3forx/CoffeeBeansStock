package handlers

import (
	"errors"
	"net/http"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/api"
	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
)

// handleDomainError maps domain errors to HTTP responses.
// notFoundMsg customizes the 404 message per resource type.
func handleDomainError(w http.ResponseWriter, err error, notFoundMsg string) {
	var validationErrs domain.ValidationErrors
	var validationErr *domain.ValidationError
	switch {
	case errors.As(err, &validationErrs):
		details := make([]api.FieldError, len(validationErrs))
		for i, ve := range validationErrs {
			details[i] = api.FieldError{Field: ve.Field, Message: ve.Message}
		}
		api.WriteValidationError(w, details)
	case errors.As(err, &validationErr):
		api.WriteValidationError(w, []api.FieldError{
			{Field: validationErr.Field, Message: validationErr.Message},
		})
	case errors.Is(err, domain.ErrNotFound):
		api.WriteError(w, http.StatusNotFound, "NOT_FOUND", notFoundMsg)
	case errors.Is(err, domain.ErrForbidden):
		api.WriteError(w, http.StatusForbidden, "FORBIDDEN", "このリソースにアクセスする権限がありません")
	case errors.Is(err, domain.ErrInsufficientStock):
		api.WriteError(w, http.StatusConflict, "INSUFFICIENT_STOCK", "在庫が不足しています")
	case errors.Is(err, domain.ErrInvalidToken), errors.Is(err, domain.ErrExpiredToken):
		api.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", notFoundMsg)
	default:
		api.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "サーバーエラーが発生しました")
	}
}
