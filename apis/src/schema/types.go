package schema

import "github.com/dgrijalva/jwt-go"

type Jwtclaims struct {
	Claims jwt.StandardClaims
}
