/*
Package auth contains the jwt authentication service middleware for the client facing
API. The decision was made to use stateless jwt, because I'm a lazy mofo.
*/
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// CreateNewUser adds a new user and returns its api key
// If a jwt.Claim is passed through, CreateUser will use that.
func CreateUser(username string, role Role, privateKey string, claims ...jwt.Claims) (string, error) {
	var babelClaims jwt.Claims

	if len(claims) > 0 {
		babelClaims = claims[0]
	} else {
		// check to see if username exists in the database and retrieve that
		// TODO

		// otherwise, generate a new claim for the user`
		babelClaims = jwt.MapClaims{
			"jti":  username,
			"role": role,
			"iat":  time.Now().Unix(),
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, babelClaims)
	return token.SignedString([]byte(privateKey))
}

// DeleteUser removes the user from the database. Attempting to delete a user that does not
// exist will return an error.
func DeleteUser(username string) error {
	return nil
}

// CreateProject creates a new project. Adding an already existing project will return an error.
func CreateProject(project_name string, email ...string) error {
	var varString string
	if len(email) == 0 {
		varString = "NULL"
	} else {
		varString = email[0]
	}
	fmt.Println(varString)
	// query := `INSERT INTO babel.projects (project_name, email)
	// VALUES`
	return nil
}

// DeleteProject deletes an existing project. Attempting to delete a project that does not exist will return an error.
func DeleteProject() error {

	return nil
}

// AddProjectAccess adds write/update access to a specific username to specific project names.
func AddProjectToUser(username string, project_names string) error {
	return nil
}

// RemoveProjectFromUser removes write/update access to a specific username to specific project names.
func RemoveProjectFromUser(username string, project_name string) error {
	return nil
}

// RetrieveAPIKey retrieves an api key for an existing user
func RetrieveAPIKey(username string) (string, error) {
	return "", nil
}
