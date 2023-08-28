package schema

import (
	"github.com/dgrijalva/jwt-go"
)

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

type CreateResturant struct {
	Name          string `json:"name"`
	Address       ResturantAddress
	Image         []ResturantImage
	ResturantType string `json:"resturant_type"`
}
type ResturantImage struct {
	Name string `json:"resturant_image"`
	Url  string `json:"resturant_image_url"`
}
type ResturantAddress struct {
	Line1     string `json:"address_line_1"`
	Line2     string `json:"address_line_2"`
	State     string `json:"state"`
	City      string `json:"city"`
	ZipCode   string `json:"zipcode"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}
