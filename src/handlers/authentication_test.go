package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/codemunsta/risevest-test/src/db"
	"github.com/codemunsta/risevest-test/src/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewTestDB() *gorm.DB {
	fmt.Print("creating")
	dsn := "host=test_db user=test password=test dbname=test_database port=5432 sslmode=disable TimeZone=Africa/Lagos"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal("Failed to setup or connect to test database. \n", err)
	}
	fmt.Print("migrating")
	db.AutoMigrate(&models.User{}, &models.Admin{}, &models.Folder{}, &models.File{}, &models.FileDelete{})

	return db
}

func TestRegisterHandler(t *testing.T) {

	fmt.Print("starting")
	testDB := NewTestDB()

	db.Database.DB = testDB
	defer func() {
		db.Database.DB = nil
		testDB.Migrator().DropTable(&models.User{}, &models.Admin{}, &models.Folder{}, &models.File{}, &models.FileDelete{})
	}()

	// Set up the request
	newUser := models.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "password",
	}
	requestBody, _ := json.Marshal(newUser)
	request := httptest.NewRequest("POST", "/register", strings.NewReader(string(requestBody)))
	writer := httptest.NewRecorder()

	// Set up the expected Create method call
	Register(writer, request)
	assert.Equal(t, http.StatusCreated, writer.Code)
	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, "User registered successfully", response["message"])
}
