package sendowl

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"

	"golang.org/x/net/context"
)

// UpdateProductRequest is an API representation of an update product request.
type UpdateProductRequest struct {
	ProductID   `json:"-"`  // ProductID of the product to update.
	Name        string      // Name of the product.
	Type        ProductType // Type of the product.
	Price       Price       // Price of the product.
	Attachment  io.Reader
	Filename    string // Filename of the file to attach.
	PDFStamping bool
}

func (r *UpdateProductRequest) body() (io.Reader, string, error) {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	w.WriteField("product[name]", r.Name)
	w.WriteField("product[product_type]", string(r.Type))
	w.WriteField("product[price]", r.Price.String())
	if r.PDFStamping {
		w.WriteField("product[pdf_stamping]", "true")
	}
	part, err := w.CreateFormFile("product[attachment]", r.Filename)
	if err != nil {
		return nil, "", err
	}
	if _, err := io.Copy(part, r.Attachment); err != nil {
		return nil, "", err
	}
	if err := w.Close(); err != nil {
		return nil, "", err
	}
	return buf, w.FormDataContentType(), nil
}

// UpdateProductResponse is an API representation of an update product response.
type UpdateProductResponse struct{}

// UpdateProduct uses req to update a product, returning a UpdateProductResponse
// and non-nil error if there was a problem.
func (c Client) UpdateProduct(ctx context.Context, req UpdateProductRequest) (*UpdateProductResponse, error) {
	body, ct, err := req.body()
	if err != nil {
		return nil, err
	}
	r, err := c.newRequest("PUT", fmt.Sprintf("./products/%s", req.ProductID), body)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", ct)
	var resp UpdateProductResponse
	if err := c.do(ctx, r, &resp); err != nil {
		if isEOFErr(err) {
			return &resp, nil
		}
		return nil, err
	}
	return &resp, nil
}

func isEOFErr(err error) bool {
	return err.Error() == "EOF"
}
