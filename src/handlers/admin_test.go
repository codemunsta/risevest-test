package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codemunsta/risevest-test/src/db"
	"github.com/codemunsta/risevest-test/src/models"
	"github.com/stretchr/testify/assert"
)

// var DBInstance *sql.DB

func TestAdminGetAllFiles(t *testing.T) {

	fmt.Print("starting")
	testDB := NewTestDB()

	db.Database.DB = testDB
	defer func() {
		// Clean up test database resources
		db.Database.DB = nil
		testDB.Migrator().DropTable(&models.User{}, &models.Admin{}, &models.Folder{}, &models.File{}, &models.FileDelete{})
	}()

	req, err := http.NewRequest("GET", "/api/admin/get/files", nil)
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(AdminGetAllFiles)

	handler.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Here are your files", response["message"])
}
