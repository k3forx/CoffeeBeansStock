package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type loggerKey struct{}

// RequestLogger returns middleware that logs each HTTP request with
// method, path, status, duration, request_id, and remote_ip.
// It also stores a request-scoped *slog.Logger (with request_id) in the context.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := chimiddleware.GetReqID(r.Context())

		logger := slog.Default().With("request_id", reqID)
		ctx := context.WithValue(r.Context(), loggerKey{}, logger)

		ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r.WithContext(ctx))

		duration := time.Since(start)
		status := ww.Status()

		attrs := []slog.Attr{
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", status),
			slog.Float64("duration_ms", float64(duration.Microseconds())/1000.0),
			slog.String("request_id", reqID),
			slog.String("remote_ip", r.RemoteAddr),
		}

		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}

		slog.LogAttrs(r.Context(), level, "request completed", attrs...)
	})
}

// LoggerFromContext returns the request-scoped logger stored by RequestLogger middleware.
// Falls back to slog.Default() if not set.
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
