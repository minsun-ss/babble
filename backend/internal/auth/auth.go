/*
Package auth contains the jwt authentication service middleware for the client facing
API. The decision was made to use stateless jwt, because I'm a lazy mofo.
*/
package auth

import (
	"babel/backend/internal/models"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

func (r Role) String() string {
	switch r {
	case RoleUser:
		return "user"
	case RoleAdmin:
		return "admin"
	default:
		return "unknown"
	}
}

func STRole(s string) (Role, error) {
	role := Role(s)
	switch role {
	case RoleUser, RoleAdmin:
		return role, nil
	default:
		return "", fmt.Errorf("invalid role: %s", s)
	}
}

// userExists checks to see if the user already exists
func userExists(db *gorm.DB, username string) bool {
	var dbUserResult []models.DBUsername

	db.Raw(`SELECT username from babel.users
		WHERE username=@username`,
		sql.Named("username", username),
	).Scan(&dbUserResult)

	if len(dbUserResult) > 0 {
		return true
	}
	return false
}

// projectExists checks the database to see if the project name already exists
// in the database
func projectExists(db *gorm.DB, project_name string) bool {
	var dbProjectResult []models.DBProjectName

	db.Raw(`SELECT project_name from babel.projects
		WHERE project_name=@project`,
		sql.Named("project", project_name),
	).Scan(&dbProjectResult)

	if len(dbProjectResult) > 0 {
		return true
	}
	return false
}

func addUser(db *gorm.DB, username string, role string, iat int64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		user := []models.DBUserInsert{
			{Username: username, Role: role, Iat: iat},
		}

		result := db.Create(&user)

		if result.Error != nil {
			return result.Error
		}
		return nil
	})
}

// CreateNewUser adds a new user and returns its api key
// If a jwt.Claim is passed through, CreateUser will use that.
func CreateUser(db *gorm.DB, private_key string, username string, role Role, claims ...jwt.Claims) (string, error) {
	var babelClaims jwt.Claims

	if db == nil {
		slog.Error("Why is the db empty")
	}
	if userExists(db, username) {
		return "", fmt.Errorf("user %s already exists, user will not be created", username)
	}

	iat := time.Now().Unix()
	err := addUser(db, username, role.String(), iat)

	if err != nil {
		return "", fmt.Errorf("error in inserting user into database, %v", err)
	}

	if len(claims) > 0 {
		babelClaims = claims[0]
	} else {
		babelClaims = jwt.MapClaims{
			"jti":  username,
			"role": role.String(),
			"iat":  iat,
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, babelClaims)
	return token.SignedString([]byte(private_key))
}

// removeUser removes a specified user from the database. This will cause existing
// privileges to specifi users to disappear.
func removeUser(db *gorm.DB, username string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		result := db.Where("username = ?", username).Delete(&models.DBUsername{})

		if result.Error != nil {
			return result.Error
		}
		return nil
	})
}

// DeleteUser removes the user from the database. Attempting to delete a user that // does not exist will return an error.
func DeleteUser(db *gorm.DB, username string) error {
	inDatabase := userExists(db, username)
	if !inDatabase {
		return fmt.Errorf("Error attempting to delete a user that does not exist")
	}

	err := removeUser(db, username)
	if err != nil {
		return fmt.Errorf("Error attempting to delete a user: %v", err)
	}
	return nil
}

func addProject(db *gorm.DB, project_name string, email []string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var projectInsert models.DBProjectInsert

		if len(email) > 0 {
			projectInsert = models.DBProjectInsert{
				ProjectName: project_name, Email: &email[0],
			}
		} else {
			projectInsert = models.DBProjectInsert{
				ProjectName: project_name, Email: nil,
			}
		}

		result := db.Create(&projectInsert)

		if result.Error != nil {
			return result.Error
		}
		return nil
	})
}

// CreateProject creates a new project. Adding an already existing project will return an error.
func CreateProject(db *gorm.DB, project_name string, email ...string) error {
	inDatabase := projectExists(db, project_name)
	if inDatabase {
		return fmt.Errorf("attempting to create a project but it already exists")
	}

	err := addProject(db, project_name, email)
	if err != nil {
		return fmt.Errorf("failure in adding the project to the database")
	}

	return nil
}

func removeProject(db *gorm.DB, project_name string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		result := db.Where("project_name = ?", project_name).Delete(&models.DBProjectName{})

		if result.Error != nil {
			return result.Error
		}
		return nil
	})
}

// DeleteProject deletes an existing project. Attempting to delete a project that does not exist will return an error.
func DeleteProject(db *gorm.DB, project_name string) error {
	inDatabase := projectExists(db, project_name)
	if !inDatabase {
		return fmt.Errorf("attempting to delete a project that does not exist")
	}

	err := removeProject(db, project_name)
	if err != nil {
		return fmt.Errorf("failure in removing project from the database")
	}
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

// func retrieveUserInfo(username string) jwt.MapClaims {

// }

// RetrieveAPIKey retrieves an api key for an existing user
func RetrieveAPIKey(db *gorm.DB, private_key string, username string) (string, error) {
	if userExists(db, username) {
		return "", fmt.Errorf("user %s already exists, user will not be created", username)
	} else {
		// TO DO fetch the key
		return "", nil

		// babelClaims = jwt.MapClaims{
		// 	"jti":  username,
		// 	"role": role,
		// 	"iat":  time.Now().Unix(),
		// }
		// token := jwt.NewWithClaims(jwt.SigningMethodHS256, babelClaims)
		// return token.SignedString([]byte(private_key))
	}
}
