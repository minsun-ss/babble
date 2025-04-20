package auth

import (
	"babel/backend/internal/testharness"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"gotest.tools/v3/assert"
)

var testdb *gorm.DB

func TestMain(m *testing.M) {
	var cleanup func()
	testdb, cleanup = testharness.SetupTestDB(m)

	if testdb == nil {
		slog.Error("There is an empty db in test main")
	}
	testharness.ResetDBData(testdb)
	if testdb == nil {
		slog.Error("There is an empty db in test main after reset data")
	}

	code := m.Run()

	cleanup()
	os.Exit(code)
}

func TestCreateKey(t *testing.T) {
	if testdb == nil {
		slog.Error("There is an empty db in test create key")
	}

	role := RoleUser
	username := "fakeuser"

	fake_claim := jwt.MapClaims{
		"username": username,
		"role":     role.String(),
		"iat":      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	private_key := "private_key"

	key, err := CreateUser(testdb, private_key, username, role, fake_claim)
	if err != nil {
		t.Errorf("Failed during the create user: %v", err)
		return
	}

	// validate that the user exists in the database
	inDatabase := userExists(testdb, username)
	assert.Equal(t, inDatabase, true, "User should now be in the database")

	// then check to see that the key is valid for this user
	generatedKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3MzU2ODk2MDAsInJvbGUiOiJ1c2VyIiwidXNlcm5hbWUiOiJmYWtldXNlciJ9.FntzGX0ctHMmV-q_aWqIUMxI742cjl6TvhWLkPhaV_Y"
	assert.Equal(t, key, generatedKey, "Creation of a static claim with a static secret key should generate the exact same signed key.")
}

func TestAddRemoveUser(t *testing.T) {
	testharness.ResetDBData(testdb)

	if testdb == nil {
		slog.Error("There is an empty db in test create key")
	}

	role := RoleUser
	username := "fakeusertest"

	fake_claim := jwt.MapClaims{
		"username": username,
		"role":     role.String(),
		"iat":      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	private_key := "private_key"

	_, err := CreateUser(testdb, private_key, username, role, fake_claim)
	if err != nil {
		t.Errorf("Failed during the create user: %v", err)
		return
	}

	// validate that the user exists in the database
	t.Log("validating now hiding that the fake user exists in the database...")
	inDatabase := userExists(testdb, username)
	assert.Equal(t, inDatabase, true, "User should now be in the database")

	// now remove the name from the database
	err = DeleteUser(testdb, username)
	if err != nil {
		t.Errorf("Failed to delet ethe user: %v", err)
		return
	}

	t.Log("validating now hiding that the fake user no longer exists in the database...")
	inDatabase = userExists(testdb, username)
	assert.Equal(t, inDatabase, false, "User should no longer be in the database")
}
