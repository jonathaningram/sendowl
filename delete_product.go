package sendowl

import (
	"fmt"

	"golang.org/x/net/context"
)

// DeleteProductRequest is an API representation of a delete product request.
type DeleteProductRequest struct {
	ProductID `json:"-"` // ProductID of the product to delete.
}

// DeleteProductResponse is an API representation of a delete product response.
type DeleteProductResponse struct{}

// DeleteProduct uses req to delete a product, returning a DeleteProductResponse
// and non-nil error if there was a problem.
func (c Client) DeleteProduct(ctx context.Context, req DeleteProductRequest) (*DeleteProductResponse, error) {
	r, err := c.newRequest("DELETE", fmt.Sprintf("./products/%s", req.ProductID), nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	r.Header.Set("Accept", "application/json")
	var resp DeleteProductResponse
	if err := c.do(ctx, r, &resp); err != nil {
		if isEOFErr(err) {
			return &resp, nil
		}
		return nil, err
	}
	return &resp, nil
}
