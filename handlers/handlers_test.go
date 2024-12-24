package handlers

import (
	"babel/models"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
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

	t.Log("Spinning up db container...")
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

	// create the docs table
	_, err = connPool.Exec(`
		CREATE TABLE IF NOT EXISTS docs (
		  name varchar(50) NOT NULL,
		  description varchar(50) DEFAULT NULL,
		  hidden tinyint(1) DEFAULT NULL,
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
		t.Fatalf("failed to create docs table: %v", err)
	}

	cleanup := func() {
		connPool.Close()
		container.Terminate(ctx)
	}

	return db, cleanup
}

// tests that the menu generated for the index is valid
func TestIndexMenu(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	connPool, err := db.DB()
	if err != nil {
		t.Fatalf("error in getting underlying database: %v", err)
	}

	t.Log("inserting test data into the database...")
	fileData, err := os.ReadFile("../test/output.zip")
	if err != nil {
		t.Fatalf("error in opening zip file: %v", err)
	}

	_, err = connPool.Exec(`
		INSERT INTO babel.docs (name, description, hidden)
		VALUES ('test1', 'testing library 1', 0),
		('test2', 'testing library 2', 1);
	`)
	if err != nil {
		t.Fatalf("failed to test values into docstable: %v", err)
	}

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

	results := generateMenuFields(db)
	assert.Equal(t, len(results), 1, "Menu should only return 1 item")
	menuItem := results[0]
	assert.Equal(t, menuItem.Title, "test1", "Menu should only return test1 library")
	assert.Equal(t, menuItem.Link, "/docs/test1", "Menu's link should be /docs/test1")
	assert.Equal(t, len(menuItem.Children), 2, "Menu should only have 2 dropdown links")
	assert.Equal(t, menuItem.MoreInfo, "/info/test1", "Menu info should be info/test1")

	t.Log("validating now hiding test1 generates no menu...")
	_, err = connPool.Exec(`
		UPDATE babel.docs
		SET hidden=1
		WHERE name="test1"
	`)
	if err != nil {
		t.Fatalf("failed to update values in docs: %v", err)
	}
	results = generateMenuFields(db)
	assert.Equal(t, len(results), 0, "There should be no menu now.")
}

func HandleLibraryPage(res http.ResponseWriter, req *http.Request) {
	data := models.LibraryData{
		Library:     "traderpythonlib",
		Description: "A trading library",
		Links: []models.LibraryLink{
			{Version: "1.2.3", Link: "/docs"},
			{Version: "1.2.1", Link: "/docs"},
		},
	}
	page := template.Must(template.ParseFiles("templates/library.html"))
	page.ExecuteTemplate(res, "library", data)
}

func TestLibraryMenu(t *testing.T) {

}

// func ServeZipFile(res http.ResponseWriter, req *http.Request) {
// 	check := html.EscapeString(req.URL.Path)
// 	fmt.Printf(check)
// 	filename := "docs/output.zip"
// 	if filename == "" {
// 		fmt.Println("where's my file?")
// 	}
// 	zr, err := zip.OpenReader(filename)
// 	if err != nil {
// 		fmt.Println("Shit my file")
// 	}

// 	defer zr.Close()
// 	prefix := "/docs/"

// 	// http.FileServer(http.FS(zr)).ServeHTTP(res, req)
// 	stripped := http.StripPrefix(prefix, http.FileServer(http.FS(zr)))
// 	stripped.ServeHTTP(res, req)
// }
