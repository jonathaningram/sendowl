package sendowl

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Price float64

// UnmarshalJSON implements json.Unmarshaler.
func (p *Price) UnmarshalJSON(data []byte) error {
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return fmt.Errorf("Price should either be a float64 or a string, got %T: %v", data, data)
		}
		f, err = strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("failed to parse Price string %q as a float: %s", s, err)
		}
	}
	*p = PriceFromFloat64(f)
	return nil
}

func (p Price) String() string {
	return strconv.FormatFloat(float64(p), 'f', -1, 64)
}

func PriceFromFloat64(f float64) Price {
	return Price(f)
}
