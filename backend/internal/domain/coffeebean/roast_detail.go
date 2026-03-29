package coffeebean

import domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"

type RoastDetail struct {
	value string
}

const roastMediumValue = "medium"

var (
	RoastDetailLight    = RoastDetail{"light"}
	RoastDetailCinnamon = RoastDetail{"cinnamon"}
	RoastDetailMedium   = RoastDetail{roastMediumValue}
	RoastDetailHigh     = RoastDetail{"high"}
	RoastDetailCity     = RoastDetail{"city"}
	RoastDetailFullCity = RoastDetail{"full_city"}
	RoastDetailFrench   = RoastDetail{"french"}
	RoastDetailItalian  = RoastDetail{"italian"}
)

func NewRoastDetail(v string) (RoastDetail, error) {
	switch v {
	case "light", "cinnamon", roastMediumValue, "high",
		"city", "full_city", "french", "italian":
		return RoastDetail{value: v}, nil
	default:
		return RoastDetail{}, &domain.ValidationError{
			Field:   "roast_detail",
			Message: "焙煎度の詳細は light/cinnamon/medium/high/city/full_city/french/italian のいずれかを指定してください",
		}
	}
}

func ReconstructRoastDetail(v string) RoastDetail { return RoastDetail{value: v} }

func (r RoastDetail) String() string { return r.value }
