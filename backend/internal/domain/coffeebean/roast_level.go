package coffeebean

import domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"

type RoastLevel struct {
	value string
}

var (
	RoastShallow    = RoastLevel{"shallow"}
	RoastMedium     = RoastLevel{roastMediumValue}
	RoastMediumDeep = RoastLevel{"medium_deep"}
	RoastDeep       = RoastLevel{"deep"}
)

func NewRoastLevel(v string) (RoastLevel, error) {
	switch v {
	case "shallow", roastMediumValue, "medium_deep", "deep":
		return RoastLevel{value: v}, nil
	default:
		return RoastLevel{}, &domain.ValidationError{
			Field:   "roast_level",
			Message: "焙煎度は shallow/medium/medium_deep/deep のいずれかを指定してください",
		}
	}
}

func ReconstructRoastLevel(v string) RoastLevel { return RoastLevel{value: v} }

func (r RoastLevel) String() string { return r.value }

// ValidDetailFor returns true if the given RoastDetail is consistent with this roast level.
func (r RoastLevel) ValidDetailFor(detail RoastDetail) bool {
	switch r.value {
	case "shallow":
		return detail.value == "light" || detail.value == "cinnamon"
	case roastMediumValue:
		return detail.value == roastMediumValue || detail.value == "high"
	case "medium_deep":
		return detail.value == "city" || detail.value == "full_city"
	case "deep":
		return detail.value == "french" || detail.value == "italian"
	default:
		return false
	}
}
