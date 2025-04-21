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

func accessExists(db *gorm.DB, username string, project_name string) bool {
	var dbProjectAccessResult []models.DBUserAccess

	db.Raw(`SELECT username, project_name from babel.user_access
		WHERE username=@username AND project_name=@project_name`,
		sql.Named("username", username),
		sql.Named("project_name", project_name)).Scan(&dbProjectAccessResult)

	if len(dbProjectAccessResult) > 0 {
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
func CreateUser(db *gorm.DB, private_key string, username string, role Role, claims ...jwt.MapClaims) (string, error) {
	var babelClaims jwt.MapClaims
	var usernameInput string
	var roleValue string
	var iatValue int64
	var ok bool

	if len(claims) > 0 {
		babelClaims = claims[0]
		iatValue, ok = babelClaims["iat"].(int64)
		if !ok {
			fmt.Errorf("Something happened in fetching iat")
		}
		roleValue, ok = babelClaims["role"].(string)
		if !ok {
			fmt.Errorf("Something happened in fetching role")
		}
		if usernameValue, ok := babelClaims["jti"].(string); ok {
			usernameInput = usernameValue
		}
	} else {
		usernameInput = username
		iatValue = time.Now().Unix()
		roleValue = role.String()
	}

	slog.Error("Sinking into database: ", "username", usernameInput, "role", roleValue, "iat", iatValue)
	if userExists(db, username) {
		return "", fmt.Errorf("user %s already exists, user will not be created", username)
	}

	err := addUser(db, usernameInput, roleValue, iatValue)

	if err != nil {
		return "", fmt.Errorf("error in inserting user into database, %v", err)
	}

	babelClaims = jwt.MapClaims{
		"jti":  usernameInput,
		"role": roleValue,
		"iat":  iatValue,
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

func addProjectToUser(db *gorm.DB, username string, project_name string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		accessInsert := models.DBUserAccess{
			Username:    username,
			ProjectName: project_name,
		}

		result := db.Create(&accessInsert)
		if result.Error != nil {
			fmt.Errorf("failure in adding access to the database")
		}

		return nil
	})
}

// AddProjectAccess adds write/update access to a specific username to specific project names.
func GrantProjectAccess(db *gorm.DB, username string, project_name string) error {
	// check to see user and project exist
	inDatabase := accessExists(db, username, project_name)
	if inDatabase {
		return fmt.Errorf("error attempting to add a privilege already granted")
	}

	err := addProjectToUser(db, username, project_name)
	if err != nil {
		return fmt.Errorf("failure in granting access to user")
	}

	inDatabase = accessExists(db, username, project_name)
	if !inDatabase {
		return fmt.Errorf("error in validating that the privilege was already granted")
	}
	return nil
}

func removeProjectFromUser(db *gorm.DB, username string, project_name string) error {
	return db.Transaction(func(db *gorm.DB) error {
		result := db.Where("username = ? AND project_name = ?", username, project_name).Delete(&models.DBUserAccess{})
		if result.Error != nil {
			return fmt.Errorf("failure in removing access to the database")
		}

		return nil
	})
}

// RevokeProjectFromUser removes write/update access to a specific username to specific project names.
func RevokeProjectAccess(db *gorm.DB, username string, project_name string) error {
	inDatabase := accessExists(db, username, project_name)
	if !inDatabase {
		return fmt.Errorf("error in removing project whose access doesn't exist")
	}

	err := removeProjectFromUser(db, username, project_name)
	if err != nil {
		return fmt.Errorf("failure in revoking project access from user")
	}

	return nil
}

func retrieveUserKey(db *gorm.DB, username string) (jwt.MapClaims, error) {
	var userResults []models.DBUserKey

	result := db.Where("username = ?", username).Find(&userResults)

	if result.Error != nil {
		return nil, fmt.Errorf("an error occurred when trying to retrieve the results")
	}

	if len(userResults) == 0 {
		return nil, fmt.Errorf("no record of this user exists in the database")
	}

	slog.Error("logging this", "username", userResults[0].Username, "role", userResults[0].Role, "iat", userResults[0].IAT)
	return jwt.MapClaims{
		"jti":  userResults[0].Username,
		"role": userResults[0].Role,
		"iat":  userResults[0].IAT,
	}, nil
}

// RetrieveAPIKey retrieves an api key for an existing user
func RetrieveAPIKey(db *gorm.DB, private_key string, username string) (string, error) {
	userClaim, err := retrieveUserKey(db, username)

	if err != nil {
		return "", fmt.Errorf("Error in retrieving API key")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	return token.SignedString([]byte(private_key))
}
