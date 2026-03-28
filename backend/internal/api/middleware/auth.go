package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/api"
	domainauth "github.com/k3forx/CoffeeBeansStock/backend/internal/domain/auth"
)

type contextKey string

const userIDKey contextKey = "userID"

// Auth returns middleware that validates JWT tokens from the Authorization header.
func Auth(tokens domainauth.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				api.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "認証トークンが必要です")
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				api.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "認証トークンの形式が不正です")
				return
			}

			claims, err := tokens.ValidateToken(parts[1])
			if err != nil {
				api.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "認証トークンが無効です")
				return
			}

			userID, err := uuid.Parse(claims.UserID)
			if err != nil {
				api.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "認証トークンが無効です")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts the authenticated user ID from the context.
func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}
