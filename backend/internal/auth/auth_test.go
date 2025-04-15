package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gotest.tools/v3/assert"
)

func TestCreateKey(t *testing.T) {
	project_teams := []string{"test"}

	fake_claim := jwt.MapClaims{
		"jti":           "fake team",
		"project_teams": project_teams,
		"iat":           time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	private_key := "private_key"

	key, _ := GenerateAPIKey(project_teams, private_key, fake_claim)
	generatedKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3MzU2ODk2MDAsImp0aSI6ImZha2UgdGVhbSIsInByb2plY3RfdGVhbXMiOlsidGVzdCJdfQ.lifJZycpOKs2rIfbHU8u9cVIPb_RLJ1YmFgns_NW28g"
	assert.Equal(t, key, generatedKey, "Creation of a static claim with a static secret key should generate the exact same signed key.")
}
