package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api"
)

// RateLimit returns middleware that limits requests per IP address.
func RateLimit(requestLimit int, windowLength time.Duration) func(http.Handler) http.Handler {
	return httprate.Limit(
		requestLimit,
		windowLength,
		httprate.WithKeyFuncs(httprate.KeyByIP),
		httprate.WithLimitHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			api.WriteError(w, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "リクエスト数の上限に達しました。しばらく経ってから再試行してください")
		})),
	)
}
