package coffeebean

import domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"

type RoastLevel struct {
	value string
}

var (
	RoastLight    = RoastLevel{"light"}
	RoastCinnamon = RoastLevel{"cinnamon"}
	RoastMedium   = RoastLevel{"medium"}
	RoastHigh     = RoastLevel{"high"}
	RoastCity     = RoastLevel{"city"}
	RoastFullCity = RoastLevel{"full_city"}
	RoastFrench   = RoastLevel{"french"}
	RoastItalian  = RoastLevel{"italian"}
)

func NewRoastLevel(v string) (RoastLevel, error) {
	switch v {
	case "light", "cinnamon", "medium", "high",
		"city", "full_city", "french", "italian":
		return RoastLevel{value: v}, nil
	default:
		return RoastLevel{}, &domain.ValidationError{
			Field:   "roast_level",
			Message: "焙煎度は light/cinnamon/medium/high/city/full_city/french/italian のいずれかを指定してください",
		}
	}
}

func ReconstructRoastLevel(v string) RoastLevel { return RoastLevel{value: v} }

func (r RoastLevel) String() string { return r.value }
