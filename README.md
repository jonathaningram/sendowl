# Sendowl API Go Client Library

Go client library for https://www.sendowl.com/developers/api/introduction

> Note: implementation is incomplete. Please submit a PR for improvements.

## Documentation

For full documentation see: [https://godoc.org/github.com/jonathaningram/sendowl](https://godoc.org/github.com/jonathaningram/sendowl).

## Install

```
go get github.com/jonathaningram/sendowl
```

## Usage

```go
package main

import (
	"log"
	"os"

	"golang.org/x/net/context"

	"github.com/jonathaningram/sendowl"
)

func main() {
	key := "sendowl-key"
	secret := "sendowl-secret"

	client := sendowl.New(*key, *secret)
	ctx := context.Background()

	filename := "file.pdf"
	name := filename
	t := sendowl.Digital
	price := 9.99

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	resp, err := client.CreateProduct(ctx, sendowl.CreateProductRequest{
		Name:       name,
		Price:      sendowl.PriceFromFloat64(price),
		Type:       t,
		Attachment: f,
		Filename:   filename,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v", resp)
}
```

To log Sendowl API requests, use `WithLogger`:

```go
client := sendowl.New(...).WithLogger(log.New(os.Stderr, "", log.LstdFlags))
```
