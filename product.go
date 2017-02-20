package sendowl

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type ProductType string

const (
	Digital ProductType = "digital"
)

type ProductID int

// UnmarshalJSON implements the json.Unmarshaler interface.
func (id *ProductID) UnmarshalJSON(data []byte) error {
	var i int64
	if err := json.Unmarshal(data, &i); err != nil {
		return fmt.Errorf("ProductID should be an int64, got %T: %v", data, data)
	}
	*id = ProductIDFromInt(i)
	return nil
}

func (id ProductID) String() string {
	return strconv.Itoa(int(id))
}

func (id ProductID) Int() int64 {
	return int64(id)
}

func ProductIDFromString(s string) ProductID {
	i, _ := strconv.Atoi(s)
	return ProductID(i)
}

func ProductIDFromInt(i int64) ProductID {
	return ProductID(i)
}

type Product struct {
	// ID of the product.
	ID ProductID `json:"id"`
	// Name of the product.
	Name string `json:"name"`
	// Type of the product.
	Type ProductType `json:"product_type"`
	// Price of the product (in dollars and cents).
	Price Price `json:"price"`
	// InstantBuyURL for purchasing the product.
	InstantBuyURL string `json:"instant_buy_url"`
	PDFStamping   bool   `json:"pdf_stamping"`
	// SelfHostedURL is the url of the file to be issued at download (only
	// useable when the product is self hosted).
	SelfHostedURL string `json:"self_hosted_url"`
}
