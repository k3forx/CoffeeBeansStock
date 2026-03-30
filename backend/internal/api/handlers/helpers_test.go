package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseQueryInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		query      string
		key        string
		defaultVal int
		min        int
		max        int
		want       int
	}{
		{
			name:       "empty string returns default",
			query:      "",
			key:        "limit",
			defaultVal: 20,
			min:        1,
			max:        100,
			want:       20,
		},
		{
			name:       "valid value within range",
			query:      "limit=50",
			key:        "limit",
			defaultVal: 20,
			min:        1,
			max:        100,
			want:       50,
		},
		{
			name:       "value below min returns min",
			query:      "limit=0",
			key:        "limit",
			defaultVal: 20,
			min:        1,
			max:        100,
			want:       1,
		},
		{
			name:       "value above max returns max",
			query:      "limit=999",
			key:        "limit",
			defaultVal: 20,
			min:        1,
			max:        100,
			want:       100,
		},
		{
			name:       "non-numeric string returns default",
			query:      "limit=abc",
			key:        "limit",
			defaultVal: 20,
			min:        1,
			max:        100,
			want:       20,
		},
		{
			name:       "negative value with min=0 returns 0",
			query:      "offset=-5",
			key:        "offset",
			defaultVal: 0,
			min:        0,
			max:        1000,
			want:       0,
		},
		{
			name:       "value equal to min returns min",
			query:      "limit=1",
			key:        "limit",
			defaultVal: 20,
			min:        1,
			max:        100,
			want:       1,
		},
		{
			name:       "value equal to max returns max",
			query:      "limit=100",
			key:        "limit",
			defaultVal: 20,
			min:        1,
			max:        100,
			want:       100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			url := "/test"
			if tt.query != "" {
				url += "?" + tt.query
			}
			r := httptest.NewRequest(http.MethodGet, url, nil)

			got := parseQueryInt(r, tt.key, tt.defaultVal, tt.min, tt.max)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("parseQueryInt() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
