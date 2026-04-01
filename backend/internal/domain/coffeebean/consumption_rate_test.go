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
	AlertLevel       coffeebean.AlertLevel
}

func toResult(cr coffeebean.ConsumptionRate) consumptionRateResult {
	return consumptionRateResult{
		RemainingCups:    cr.RemainingCups(),
		RemainingDays:    cr.RemainingDays(),
		DailyConsumption: cr.DailyConsumption(),
		WeeklyTotal:      cr.WeeklyTotal(),
		AlertLevel:       cr.AlertLevel(),
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
				AlertLevel:       coffeebean.AlertLevelNone,
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
				AlertLevel:       coffeebean.AlertLevelDanger,
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
				AlertLevel:       coffeebean.AlertLevelNone,
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
				AlertLevel:       coffeebean.AlertLevelNone,
			},
		},
		"残日数3日以下の場合_アラートはdanger": {
			currentStock: 45,
			gramsPerCup:  15,
			weeklyUsage:  ptr(int32(105)),
			want: consumptionRateResult{
				RemainingCups:    3,
				RemainingDays:    ptr(int32(3)),
				DailyConsumption: ptr(15.0),
				WeeklyTotal:      ptr(int32(105)),
				AlertLevel:       coffeebean.AlertLevelDanger,
			},
		},
		"残日数4日の場合_アラートはwarning": {
			currentStock: 60,
			gramsPerCup:  15,
			weeklyUsage:  ptr(int32(105)),
			want: consumptionRateResult{
				RemainingCups:    4,
				RemainingDays:    ptr(int32(4)),
				DailyConsumption: ptr(15.0),
				WeeklyTotal:      ptr(int32(105)),
				AlertLevel:       coffeebean.AlertLevelWarning,
			},
		},
		"残日数7日の場合_アラートはwarning": {
			currentStock: 105,
			gramsPerCup:  15,
			weeklyUsage:  ptr(int32(105)),
			want: consumptionRateResult{
				RemainingCups:    7,
				RemainingDays:    ptr(int32(7)),
				DailyConsumption: ptr(15.0),
				WeeklyTotal:      ptr(int32(105)),
				AlertLevel:       coffeebean.AlertLevelWarning,
			},
		},
		"残日数8日の場合_アラートなし": {
			currentStock: 120,
			gramsPerCup:  15,
			weeklyUsage:  ptr(int32(105)),
			want: consumptionRateResult{
				RemainingCups:    8,
				RemainingDays:    ptr(int32(8)),
				DailyConsumption: ptr(15.0),
				WeeklyTotal:      ptr(int32(105)),
				AlertLevel:       coffeebean.AlertLevelNone,
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
