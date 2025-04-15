/*
Package auth contains the jwt authentication service middleware for the client facing
API. The decision was made to use stateless jwt, because I'm a lazy mofo.
*/
package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// generateAPIKey generates an api key with the privileges to review certain packages
func GenerateAPIKey(project_teams []string, privateKey string, claims ...jwt.Claims) (string, error) {
	// set up claims for the key. Some general concepts
	// keys are forever (no exp)
	// we keep revocation lists in the database, however
	var babelClaims jwt.Claims
	if len(claims) > 0 {
		babelClaims = claims[0]
	} else {
		babelClaims = jwt.MapClaims{
			"jti":         uuid.New().String(),
			"permissions": project_teams,
			"iat":         time.Now().Unix(),
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, babelClaims)
	return token.SignedString([]byte(privateKey))
}
