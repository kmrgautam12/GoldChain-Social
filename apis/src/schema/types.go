package schema

import "github.com/dgrijalva/jwt-go"

type Jwtclaims struct {
	Claims jwt.StandardClaims
}

type AddItemToCart struct {
	ID      string `json:"id"`
	Payload Payload
}
type Payload struct {
	ResturantId   string  `json:"resturant_id"`
	ResturantName string  `json:"resturant_name"`
	ItemId        string  `json:"item_id"`
	ItemName      string  `json:"item_name"`
	Price         float64 `json:"price"`
}
