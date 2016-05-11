package sendowl

import (
	"fmt"

	"golang.org/x/net/context"
)

// GetProductRequest is an API representation of a get product request.
type GetProductRequest struct {
	ID ProductID
}

// GetProductResponse is an API representation of a get product response.
type GetProductResponse struct {
	Product `json:"product"`
}

// GetProduct uses req to get a product, returning a GetProductResponse and
// non-nil error if there was a problem.
func (c Client) GetProduct(ctx context.Context, req GetProductRequest) (*GetProductResponse, error) {
	u := fmt.Sprintf("./products/%s", req.ID)
	r, err := c.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	var resp GetProductResponse
	if err := c.do(ctx, r, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
