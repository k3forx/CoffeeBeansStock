package usagehistory_test

import (
	"errors"
	"testing"

	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/usagehistory"
)

func TestNewUsageType(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   string
		want    usagehistory.UsageType
		wantErr bool
	}{
		"manualは有効":       {input: "manual", want: usagehistory.UsageTypeManual},
		"quick_buttonは有効": {input: "quick_button", want: usagehistory.UsageTypeQuickButton},
		"空文字はエラー":         {input: "", wantErr: true},
		"不正な値はエラー":        {input: "invalid", wantErr: true},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := usagehistory.NewUsageType(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
					return
				}
				var ve *domain.ValidationError
				if !errors.As(err, &ve) {
					t.Errorf("expected ValidationError, got %T: %v", err, err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("NewUsageType(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestUsageType_String(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		ut   usagehistory.UsageType
		want string
	}{
		"manual":       {ut: usagehistory.UsageTypeManual, want: "manual"},
		"quick_button": {ut: usagehistory.UsageTypeQuickButton, want: "quick_button"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := tt.ut.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReconstructUsageType(t *testing.T) {
	t.Parallel()

	got := usagehistory.ReconstructUsageType("manual")
	if got.String() != "manual" {
		t.Errorf("ReconstructUsageType(\"manual\").String() = %q, want \"manual\"", got.String())
	}
}
