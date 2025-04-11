package handlers

import (
	"archive/zip"
	"babel/backend/internal/testharness"
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// tests that the menu generated for the index is valid
func TestIndexMenu(t *testing.T) {
	db, cleanup := testharness.SetupTestDB(t)
	defer cleanup()

	connPool, err := db.DB()
	if err != nil {
		t.Fatalf("error in getting underlying database: %v", err)
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
		SET is_visible=0
		WHERE name="test1"
	`)
	if err != nil {
		t.Fatalf("failed to update values in docs: %v", err)
	}
	results = generateMenuFields(db)
	assert.Equal(t, len(results), 0, "There should be no menu now.")
}

// TestLibraryMenu tests that the /info page is correctly rendered
func TestLibraryMenu(t *testing.T) {
	db, cleanup := testharness.SetupTestDB(t)
	defer cleanup()

	results := generateLibraryInfo(db, "test1")
	assert.Equal(t, results.Library, "test1", "Library names should be test2")
	assert.Equal(t, results.Description, "testing library 1", "Library descriptio should match")
	assert.Equal(t, len(results.Links), 1, "There should only be 1 link generated for test1 library")
}

// TestDocsMenu tests that the docs are correctly rendered
func TestDocsMenu(t *testing.T) {
	db, cleanup := testharness.SetupTestDB(t)
	defer cleanup()

	t.Log("check to see that a valid zip file can be opened...")
	data, err := generateDocsData(db, "test1", "3.2.1")
	assert.Nil(t, err)
	reader := bytes.NewReader(data.DataZip)

	// Try to open as ZIP and assert no nil error
	_, err = zip.NewReader(reader, int64(len(data.DataZip)))
	assert.Nil(t, err)

	t.Log("check to see that bad version data throws an error...")
	_, err = generateDocsData(db, "test2", "e.1.2")
	assert.NotNil(t, err)
}

func TestHandleHealthCheck(t *testing.T) {
	r := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	LivenessHandler(w, r)
	assert.Equal(t, w.Result().StatusCode, 200, "Health check should return OK")
}
