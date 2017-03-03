package sendowl

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type OrderID int

// UnmarshalJSON implements the json.Unmarshaler interface.
func (id *OrderID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*id = OrderIDFromString(s)
		return nil
	}
	var i int64
	if err := json.Unmarshal(data, &i); err != nil {
		return fmt.Errorf("sendowl: OrderID should be an int64, got %T: %v", data, data)
	}
	*id = OrderIDFromInt(i)
	return nil
}

func (id OrderID) String() string {
	return strconv.Itoa(int(id))
}

func (id OrderID) Int() int64 {
	return int64(id)
}

func OrderIDFromString(s string) OrderID {
	i, _ := strconv.Atoi(s)
	return OrderID(i)
}

func OrderIDFromInt(i int64) OrderID {
	return OrderID(i)
}

type Order struct {
	ID         OrderID `json:"id"` // ID of the order.
	BuyerEmail string  `json:"buyer_email"`
	BuyerName  string  `json:"buyer_name"`
	Cart       struct {
		Items []struct {
			Item struct {
				ProductID        `json:"product_id"`
				PackageID        int        `json:"package_id"`
				Quantity         int        `json:"quantity"`
				PriceAtCheckout  Price      `json:"price_at_checkout"`
				ValidUntil       *time.Time `json:"valid_until"`
				DownloadAttempts int        `json:"download_attempts"`
			} `json:"cart_item"`

			// Note: webhook returns the data in these fields directly, not in
			// the "cart_item" field above.
			// TODO(jon): probably best to find a better way to do this. Either
			// a "webhook" order needs to be used instead, or maybe the Order
			// type can be initialized with a webhook=true bool, and then
			// implement UnmarshalJSON to decode the different JSON schemas.

			Product          Product    `json:"product"`
			Quantity         int        `json:"quantity"`
			ValidUntil       *time.Time `json:"valid_until"`
			DownloadAttempts int        `json:"download_attempts"`
		} `json:"cart_items"`
		CompletedCheckoutAt time.Time `json:"completed_checkout_at"`
		StartedCheckoutAt   time.Time `json:"started_checkout_at"`
	} `json:"cart"`
}
