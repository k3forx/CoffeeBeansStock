package coffeebean

import "math"

// ConsumptionRate holds consumption metrics for a coffee bean.
type ConsumptionRate struct {
	remainingCups    int32
	remainingDays    *int32
	dailyConsumption *float64
	weeklyTotal      *int32
}

// NewConsumptionRate calculates consumption rate from stock, grams per cup, and weekly usage.
// weeklyUsage is nil when no usage data exists.
func NewConsumptionRate(currentStock, gramsPerCup int32, weeklyUsage *int32) ConsumptionRate {
	var cups int32
	if gramsPerCup > 0 {
		cups = currentStock / gramsPerCup
	}

	if weeklyUsage == nil || *weeklyUsage <= 0 {
		return ConsumptionRate{
			remainingCups:    cups,
			remainingDays:    nil,
			dailyConsumption: nil,
			weeklyTotal:      nil,
		}
	}

	wt := *weeklyUsage
	daily := float64(wt) / 7.0
	days := int32(math.Floor(float64(currentStock) / daily))

	return ConsumptionRate{
		remainingCups:    cups,
		remainingDays:    &days,
		dailyConsumption: &daily,
		weeklyTotal:      &wt,
	}
}

func (c ConsumptionRate) RemainingCups() int32       { return c.remainingCups }
func (c ConsumptionRate) RemainingDays() *int32      { return c.remainingDays }
func (c ConsumptionRate) DailyConsumption() *float64 { return c.dailyConsumption }
func (c ConsumptionRate) WeeklyTotal() *int32        { return c.weeklyTotal }
