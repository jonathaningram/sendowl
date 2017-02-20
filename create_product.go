package sendowl

import (
	"bytes"
	"io"
	"mime/multipart"

	"golang.org/x/net/context"
)

// CreateProductRequest is an API representation of a create product request.
type CreateProductRequest struct {
	Name        string      // Name of the product.
	Type        ProductType // Type of the product.
	Price       Price       // Price of the product.
	Attachment  io.Reader
	Filename    string // Filename of the file to attach.
	PDFStamping bool
	// SelfHostedURL is the url of the file to be issued at download (only
	// useable when the product is self hosted).
	SelfHostedURL string
}

func (r *CreateProductRequest) body() (io.Reader, string, error) {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	w.WriteField("product[name]", r.Name)
	w.WriteField("product[product_type]", string(r.Type))
	w.WriteField("product[price]", r.Price.String())
	if r.PDFStamping {
		w.WriteField("product[pdf_stamping]", "true")
	}
	w.WriteField("product[self_hosted_url]", r.SelfHostedURL)
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

// CreateProductResponse is an API representation of a create product response.
type CreateProductResponse struct {
	Product Product `json:"product"`
}

// CreateProduct uses req to create a new product, returning a
// CreateProductResponse and non-nil error if there was a problem.
func (c Client) CreateProduct(ctx context.Context, req CreateProductRequest) (*CreateProductResponse, error) {
	body, ct, err := req.body()
	if err != nil {
		return nil, err
	}
	r, err := c.newRequest("POST", "./products.json", body)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", ct)
	var resp CreateProductResponse
	if err := c.do(ctx, r, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
