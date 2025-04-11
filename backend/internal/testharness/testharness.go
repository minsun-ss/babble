/* testharness package is a collection of utilities to run the tests in this package.
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
func SetupTestDB(t *testing.T) (*gorm.DB, func()) {
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

	t.Log("spinning up db container...")
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	dbPort, err := container.MappedPort(ctx, "3306")
	if err != nil {
		t.Fatalf("failed to get container port: %v", err)
	}

	t.Log("connecting to db container...")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		"myuser", "mypassword", host, dbPort.Port(), "babel")

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		t.Fatalf("no connection established test database: %v", err)
	}

	connPool, err := db.DB()
	if err != nil {
		t.Fatalf("no connection established to test database: %v", err)
	}

	connPool.SetMaxOpenConns(20)
	connPool.SetConnMaxLifetime(time.Hour)

	// create the docs and doc history table
	t.Log("creating the babel tables to test against...")
	_, err = connPool.Exec(`
		CREATE TABLE IF NOT EXISTS docs (
		  name varchar(50) NOT NULL,
		  description varchar(50) DEFAULT NULL,
		  is_visible tinyint(1) DEFAULT NULL,
		  project_team varchar(50) DEFAULT "Other",
		  last_updated_dt timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
		  PRIMARY KEY (name),
		  KEY ix_last_updated_dt (last_updated_dt)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

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
		t.Fatalf("failed to create docs and doc_history tables: %v", err)
	}

	t.Log("inserting test data into the database...")
	fileData, err := os.ReadFile("../testharness/test.zip")
	if err != nil {
		t.Fatalf("error in opening zip file: %v", err)
	}

	// populate with some fake data
	_, err = connPool.Exec(`
		INSERT INTO babel.docs (name, description, is_visible)
		VALUES ('test1', 'testing library 1', 1),
		('test2', 'testing library 2', 0);
	`)
	if err != nil {
		t.Fatalf("failed to test values into docstable: %v", err)
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
		t.Fatalf("failed to test values into docs and docs_history table: %v", err)
	}

	// generate cleanup function to call later to remove the container after testing
	cleanup := func() {
		connPool.Close()
		container.Terminate(ctx)
	}

	return db, cleanup
}
