/*
Package auth contains the jwt authentication service middleware for the client facing
API. The decision was made to use stateless jwt, because I'm a lazy mofo.
*/
package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

func generateAPIKey() (string, error) {
	claims := jwt.MapClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString([]byte(secretKey))
}
