package usagehistory

import domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"

// UsageType represents the type of usage (manual or quick_button).
type UsageType struct {
	value string
}

var (
	UsageTypeManual      = UsageType{value: "manual"}
	UsageTypeQuickButton = UsageType{value: "quick_button"}
)

// NewUsageType creates a new UsageType with validation.
func NewUsageType(s string) (UsageType, error) {
	switch s {
	case "manual":
		return UsageTypeManual, nil
	case "quick_button":
		return UsageTypeQuickButton, nil
	default:
		return UsageType{}, &domain.ValidationError{
			Field:   "usage_type",
			Message: "使用タイプは manual または quick_button を指定してください",
		}
	}
}

// ReconstructUsageType restores a UsageType from persisted data without validation.
func ReconstructUsageType(s string) UsageType { return UsageType{value: s} }

// String returns the string representation of the UsageType.
func (t UsageType) String() string { return t.value }
