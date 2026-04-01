package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/api"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api/middleware"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/apperrors"
	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
)

// handleDomainError maps domain errors to HTTP responses and logs them.
// notFoundMsg customizes the 404/401 message per resource type.
func handleDomainError(w http.ResponseWriter, r *http.Request, err error, notFoundMsg string) {
	var validationErrs domain.ValidationErrors
	var validationErr *domain.ValidationError
	switch {
	case errors.As(err, &validationErrs):
		logError(r.Context(), slog.LevelWarn, err)
		details := make([]api.FieldError, len(validationErrs))
		for i, ve := range validationErrs {
			details[i] = api.FieldError{Field: ve.Field, Message: ve.Message}
		}
		api.WriteValidationError(w, details)
	case errors.As(err, &validationErr):
		logError(r.Context(), slog.LevelWarn, err)
		api.WriteValidationError(w, []api.FieldError{
			{Field: validationErr.Field, Message: validationErr.Message},
		})
	case errors.Is(err, domain.ErrNotFound):
		logError(r.Context(), slog.LevelWarn, err)
		api.WriteError(w, http.StatusNotFound, "NOT_FOUND", notFoundMsg)
	case errors.Is(err, domain.ErrForbidden):
		logError(r.Context(), slog.LevelWarn, err)
		api.WriteError(w, http.StatusForbidden, "FORBIDDEN", "このリソースにアクセスする権限がありません")
	case errors.Is(err, domain.ErrInsufficientStock):
		logError(r.Context(), slog.LevelWarn, err)
		api.WriteError(w, http.StatusConflict, "INSUFFICIENT_STOCK", "在庫が不足しています")
	case errors.Is(err, domain.ErrInvalidToken), errors.Is(err, domain.ErrExpiredToken):
		logError(r.Context(), slog.LevelWarn, err)
		api.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", notFoundMsg)
	default:
		logError(r.Context(), slog.LevelError, err)
		api.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "サーバーエラーが発生しました")
	}
}

// logError logs an error at the specified level with optional stack trace.
func logError(ctx context.Context, level slog.Level, err error) {
	attrs := []slog.Attr{slog.String("error", err.Error())}
	if stack, ok := apperrors.GetStackTrace(err); ok {
		frames := make([]any, len(stack))
		for i, f := range stack {
			frames[i] = slog.GroupValue(
				slog.String("func", f.Function),
				slog.String("file", f.File),
				slog.Int("line", f.Line),
			)
		}
		attrs = append(attrs, slog.Any("stack_trace", frames))
	}
	middleware.LoggerFromContext(ctx).LogAttrs(ctx, level, "request error", attrs...)
}
