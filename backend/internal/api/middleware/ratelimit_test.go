package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api/middleware"
)

func TestRateLimit(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		requestLimit int
		numRequests  int
		wantStatus   int
	}{
		"under_limit": {
			requestLimit: 3,
			numRequests:  2,
			wantStatus:   http.StatusOK,
		},
		"at_limit": {
			requestLimit: 3,
			numRequests:  3,
			wantStatus:   http.StatusOK,
		},
		"over_limit": {
			requestLimit: 3,
			numRequests:  4,
			wantStatus:   http.StatusTooManyRequests,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			handler := middleware.RateLimit(tt.requestLimit, time.Minute)(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			)

			var lastRecorder *httptest.ResponseRecorder
			for i := range tt.numRequests {
				_ = i
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.RemoteAddr = "192.0.2.1:1234"
				handler.ServeHTTP(rec, req)
				lastRecorder = rec
			}

			if lastRecorder.Code != tt.wantStatus {
				t.Errorf("status code: got %d, want %d", lastRecorder.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusTooManyRequests {
				var got api.Response
				if err := json.NewDecoder(lastRecorder.Body).Decode(&got); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				want := api.Response{
					Success: false,
					Error: &api.ErrorBody{
						Code:    "RATE_LIMIT_EXCEEDED",
						Message: "リクエスト数の上限に達しました。しばらく経ってから再試行してください",
					},
				}
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("response mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
