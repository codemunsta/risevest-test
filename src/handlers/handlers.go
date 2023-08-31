package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/codemunsta/risevest-test/src/db"
	"github.com/codemunsta/risevest-test/src/models"
	"github.com/codemunsta/risevest-test/src/utils"
)

func NewApi(writer http.ResponseWriter, request *http.Request) {
	response := map[string]interface{}{
		"User":    request.RemoteAddr,
		"Message": "Pong New",
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}

func TestUserAuth(writer http.ResponseWriter, request *http.Request) {
	authToken := request.Header.Get("Authorization")
	if authToken == "" {
		http.Error(writer, "Missing authentication token", http.StatusUnauthorized)
		return
	}

	tokenString := strings.Replace(authToken, "Bearer ", "", 1)
	userID, err := utils.ParseAuthToken(tokenString)
	if err != nil {
		http.Error(writer, "Invalid authentication token", http.StatusUnauthorized)
		return
	}

	var authuser models.User

	db := db.Database.DB
	err = db.Where("ID = ?", userID).First(&authuser).Error
	if err != nil {
		http.Error(writer, "Invalid authentication token", http.StatusUnauthorized)
		return
	}

	writer.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"User name": authuser.FirstName,
	}
	json.NewEncoder(writer).Encode(response)
}
