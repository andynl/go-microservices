// Package classification of Product API
//
// Documentation for Product API
//
// Schemes : http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/andynl/go-microservices/data"
	"github.com/gorilla/mux"
)

// A list of products return in the response
// swagger:response productsResponse
type productsResponse struct {
	// All products in the system
	// in: body
	Body []data.Product
}

type KeyProduct struct{}

type Products struct {
	l *log.Logger
	v *data.Validation
}

func NewProducts(l *log.Logger, v *data.Validation) *Products {
	return &Products{l, v}
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

type ValidationError struct {
	Messages []string `json:"message"`
}

// getProductID returns the product ID from the URL
// Panics if cannot convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number
func getProductID(r *http.Request) int {
	// parse the product id from the url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}

	return id
}
