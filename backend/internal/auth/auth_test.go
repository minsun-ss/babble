package auth

import (
	"babel/backend/internal/testharness"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"gotest.tools/v3/assert"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	db, cleanup := testharness.SetupTestDB(m)

	testharness.ResetDBData(db)
	code := m.Run()

	cleanup()
	os.Exit(code)
}

func TestCreateKey(t *testing.T) {
	role := RoleUser
	username := "fakeuser"

	fake_claim := jwt.MapClaims{
		"username": username,
		"role":     role,
		"iat":      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	private_key := "private_key"

	key, _ := CreateUser(db, private_key, username, role, fake_claim)
	generatedKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3MzU2ODk2MDAsInJvbGUiOiJ1c2VyIiwidXNlcm5hbWUiOiJmYWtldXNlciJ9.FntzGX0ctHMmV-q_aWqIUMxI742cjl6TvhWLkPhaV_Y"
	assert.Equal(t, key, generatedKey, "Creation of a static claim with a static secret key should generate the exact same signed key.")
}

// func TestAdduser(t *testing.T)
