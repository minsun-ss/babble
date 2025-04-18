package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gotest.tools/v3/assert"
)

func TestCreateKey(t *testing.T) {
	role := RoleUser
	username := "fakeuser"

	fake_claim := jwt.MapClaims{
		"username": username,
		"role":     role,
		"iat":      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	private_key := "private_key"

	key, _ := CreateUser(username, role, private_key, fake_claim)
	generatedKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3MzU2ODk2MDAsInJvbGUiOiJ1c2VyIiwidXNlcm5hbWUiOiJmYWtldXNlciJ9.FntzGX0ctHMmV-q_aWqIUMxI742cjl6TvhWLkPhaV_Y"
	assert.Equal(t, key, generatedKey, "Creation of a static claim with a static secret key should generate the exact same signed key.")
}
