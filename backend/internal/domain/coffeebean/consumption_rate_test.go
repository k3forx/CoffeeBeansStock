package coffeebean_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/coffeebean"
)

func ptr[T any](v T) *T { return &v }

type consumptionRateResult struct {
	RemainingCups    int32
	RemainingDays    *int32
	DailyConsumption *float64
	WeeklyTotal      *int32
}

func toResult(cr coffeebean.ConsumptionRate) consumptionRateResult {
	return consumptionRateResult{
		RemainingCups:    cr.RemainingCups(),
		RemainingDays:    cr.RemainingDays(),
		DailyConsumption: cr.DailyConsumption(),
		WeeklyTotal:      cr.WeeklyTotal(),
	}
}

func TestNewConsumptionRate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		currentStock int32
		gramsPerCup  int32
		weeklyUsage  *int32
		want         consumptionRateResult
	}{
		"在庫150g_1杯15g_週間105g使用の場合": {
			currentStock: 150,
			gramsPerCup:  15,
			weeklyUsage:  ptr(int32(105)),
			want: consumptionRateResult{
				RemainingCups:    10,
				RemainingDays:    ptr(int32(10)),
				DailyConsumption: ptr(15.0),
				WeeklyTotal:      ptr(int32(105)),
			},
		},
		"在庫0gの場合": {
			currentStock: 0,
			gramsPerCup:  15,
			weeklyUsage:  ptr(int32(105)),
			want: consumptionRateResult{
				RemainingCups:    0,
				RemainingDays:    ptr(int32(0)),
				DailyConsumption: ptr(15.0),
				WeeklyTotal:      ptr(int32(105)),
			},
		},
		"週間使用量が0の場合_予測フィールドはnil": {
			currentStock: 150,
			gramsPerCup:  15,
			weeklyUsage:  ptr(int32(0)),
			want: consumptionRateResult{
				RemainingCups:    10,
				RemainingDays:    nil,
				DailyConsumption: nil,
				WeeklyTotal:      nil,
			},
		},
		"週間使用データがnilの場合_予測フィールドはnil": {
			currentStock: 150,
			gramsPerCup:  15,
			weeklyUsage:  nil,
			want: consumptionRateResult{
				RemainingCups:    10,
				RemainingDays:    nil,
				DailyConsumption: nil,
				WeeklyTotal:      nil,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := coffeebean.NewConsumptionRate(tt.currentStock, tt.gramsPerCup, tt.weeklyUsage)
			result := toResult(got)

			if diff := cmp.Diff(tt.want, result, cmpopts.EquateApprox(0, 0.001)); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
