/* testharness package is a collection of utilities to run the tests in this package. It is largely used to setup and teardown the mock DB and reuse.
 */
package testharness

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// SetupTestDB sets up a mariadb database container and tables to test against.
// Returns both a gorm DB connection as well as a function for its cleanup after
// completion of testing.
func SetupTestDB(m *testing.M) (*gorm.DB, func()) {
	ctx := context.Background()

	// spin up a container and wait for port to be exposed
	req := testcontainers.ContainerRequest{
		Image:        "mariadb:10.6",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MARIADB_ALLOW_EMPTY_ROOT_PASSWORD": "1",
			"MARIADB_ROOT_HOST":                 "%",
			"MARIADB_USER":                      "myuser",
			"MARIADB_PASSWORD":                  "mypassword",
			"MARIADB_DATABASE":                  "babel",
		},
		WaitingFor: wait.ForExposedPort().
			WithStartupTimeout(time.Second * 30),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		fmt.Printf("failed to start container: %v", err)
		os.Exit(1)
	}

	host, err := container.Host(ctx)
	if err != nil {
		fmt.Printf("failed to get container host: %v", err)
		os.Exit(1)
	}

	dbPort, err := container.MappedPort(ctx, "3306")
	if err != nil {
		fmt.Printf("failed to get container port: %v", err)
		os.Exit(1)
	}

	time.Sleep(2 * time.Second)
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		"myuser", "mypassword", host, dbPort.Port(), "babel")

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		fmt.Printf("no connection established test database: %v", err)
		os.Exit(1)
	}

	connPool, err := db.DB()
	if err != nil {
		fmt.Printf("no connection established to test database: %v", err)
		os.Exit(1)
	}

	connPool.SetMaxOpenConns(20)
	connPool.SetConnMaxLifetime(time.Hour)

	// create the docs and doc history table
	_, err = connPool.Exec(`
		CREATE TABLE users (
		  id bigint(20) NOT NULL AUTO_INCREMENT,
		  username varchar(50) NOT NULL,
		  iat int(11) NOT NULL,
		  last_updated_dt timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
		  role varchar(50) NOT NULL,
		  PRIMARY KEY (id),
		  UNIQUE KEY username (username)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

		CREATE TABLE docs (
		  name varchar(50) NOT NULL,
		  description varchar(50) DEFAULT NULL,
		  last_updated_dt timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
		  is_visible tinyint(4) NOT NULL DEFAULT 1,
		  id bigint(20) NOT NULL AUTO_INCREMENT,
		  project_name varchar(50) DEFAULT 'Other',
		  PRIMARY KEY (id),
		  UNIQUE KEY name (name),
		  KEY ix_last_updated_dt (last_updated_dt),
		  KEY ix_name (name),
		  KEY ix_project_name (project_name)
		) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

		CREATE TABLE doc_history (
		  id bigint(20) NOT NULL AUTO_INCREMENT,
		  name varchar(50) NOT NULL,
		  version_major int NOT NULL,
		  version_minor int NOT NULL,
		  version_patch int NOT NULL,
		  html longblob NOT NULL,
		  last_updated_dt timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
		  PRIMARY KEY (id),
		  UNIQUE KEY name (name,version_major,version_minor,version_patch),
		  KEY idx_version_patch (version_patch),
		  KEY idx_version_major (version_major),
		  KEY idx_version_minor (version_minor),
		  KEY idx_name (name)
		) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	`)
	if err != nil {
		fmt.Printf("failed to create docs and doc_history tables: %v", err)
		os.Exit(1)
	}

	// generate cleanup function to call later to remove the container after testing
	cleanup := func() {
		connPool.Close()
		container.Terminate(ctx)
	}

	return db, cleanup

}

func ResetDBData(db *gorm.DB) error {
	connPool, err := db.DB()

	if err != nil {
		return fmt.Errorf("null connPool on resetting data: %v", err)
	}

	connPool.SetMaxOpenConns(20)
	connPool.SetConnMaxLifetime(time.Hour)

	fileData, err := os.ReadFile("../testharness/test.zip")
	if err != nil {
		return fmt.Errorf("error in opening zip file: %v", err)
	}

	_, err = connPool.Exec(`
		TRUNCATE babel.users;
		TRUNCATE babel.docs;
		TRUNCATE babel.doc_history;
	`)
	if err != nil {
		return fmt.Errorf("failed to cleanup babel docs: %v", err)
	}

	// populate with some fake data
	_, err = connPool.Exec(`
		INSERT INTO babel.docs (name, description, is_visible)
		VALUES ('test1', 'testing library 1', 1),
		('test2', 'testing library 2', 0);
	`)
	if err != nil {
		return fmt.Errorf("failed to test values into docs table: %v", err)
	}

	// populate with some fake data
	_, err = connPool.Exec(`
		INSERT INTO babel.doc_history (name, version_major, version_minor, version_patch, html)
		VALUES (?, ?, ?, ?, ?)
	`,
		"test1",
		3,
		2,
		1,
		fileData,
	)
	if err != nil {
		return fmt.Errorf("failed to test value in docs_history: %v", err)
	}

	return nil
}
