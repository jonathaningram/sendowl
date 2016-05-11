package sendowl

import (
	"fmt"

	"golang.org/x/net/context"
)

// ListProductsRequest is an API representation of a list products request.
type ListProductsRequest struct {
	PerPage int
	Page    int
}

// ListProductsResponse is an API representation of a list products response.
type ListProductsResponse []listProductsResponseItem

type listProductsResponseItem struct {
	Product `json:"product"`
}

// ListProducts uses req to list a products, returning a
// ListProductsResponse and non-nil error if there was a problem.
func (c Client) ListProducts(ctx context.Context, req ListProductsRequest) (*ListProductsResponse, error) {
	u := "./products"
	if req.PerPage > 0 || req.Page > 0 {
		u = fmt.Sprintf("%s?per_page=%d&page=%d", u, req.PerPage, req.Page)
	}
	r, err := c.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	var resp ListProductsResponse
	if err := c.do(ctx, r, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
