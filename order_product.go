package sendowl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"golang.org/x/net/context"
)

// OrderProductRequest is an API representation of an order product request.
type OrderProductRequest struct {
	ProductID `json:"-"`   // ProductID of the product to order.
	Order     ProductOrder `json:"order"` // Order details.
}

// ProductOrder is an API representation of the order details for an order
// product request.
type ProductOrder struct {
	BuyerEmail       string     `json:"buyer_email"`
	BuyerName        string     `json:"buyer_name"`
	BuyerID          string     `json:"buyer_id,omitempty"`
	BuyerStatus      string     `json:"buyer_status,omitempty"`
	BuyerAddress1    string     `json:"buyer_address1,omitempty"`
	BuyerAddress2    string     `json:"buyer_address2,omitempty"`
	BuyerCity        string     `json:"buyer_city,omitempty"`
	BuyerRegion      string     `json:"buyer_region,omitempty"`
	BuyerPostCode    string     `json:"buyer_postcode,omitempty"`
	BuyerCountry     string     `json:"buyer_country,omitempty"`
	BuyerIPAddress   string     `json:"buyer_ip_address,omitempty"`
	CanMarketToBuyer bool       `json:"can_market_to_buyer"`
	Quantity         int        `json:"quantity,omitempty"`
	DispatchedAt     *time.Time `json:"dispatched_at,omitempty"`
	Tag              string     `json:"tag,omitempty"`
}

func (r *OrderProductRequest) body() (io.Reader, string, error) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(&r); err != nil {
		return nil, "", err
	}
	return buf, "application/json; charset=utf-8", nil
}

// OrderProductResponse is an API representation of an order product response.
type OrderProductResponse struct {
	Order `json:"order"`
}

// OrderProduct uses req to issue a new order, returning a
// OrderProductResponse and non-nil error if there was a problem.
func (c Client) OrderProduct(ctx context.Context, req OrderProductRequest) (*OrderProductResponse, error) {
	body, ct, err := req.body()
	if err != nil {
		return nil, err
	}
	r, err := c.newRequest("POST", fmt.Sprintf("./products/%s/issue", req.ProductID), body)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", ct)
	var resp OrderProductResponse
	if err := c.do(ctx, r, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
