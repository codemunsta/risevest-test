package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/codemunsta/risevest-test/src/db"
	"github.com/codemunsta/risevest-test/src/models"
	"github.com/codemunsta/risevest-test/src/utils"
)

func isAuthenticated(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		authToken := request.Header.Get("Authorization")
		if authToken == "" {
			NotAuthenticated(writer)
		}

		// check token in redis
		_, isFound := db.GetSession(authToken)
		if !isFound {
			tokenString := strings.Replace(authToken, "Bearer ", "", 1)
			userID, err := utils.ParseAuthToken(tokenString)
			if err != nil {
				NotAuthenticated(writer)
			}

			var authuser models.User

			database := db.Database.DB
			err = database.Where("ID = ?", userID).First(&authuser).Error
			if err != nil {
				NotAuthenticated(writer)
			}

			err = db.CreateSession(authToken, authuser)
			if err == nil {
				handler.ServeHTTP(writer, request)
			} else {
				response := map[string]interface{}{
					"message": "An Error Occured",
					"data":    nil,
				}
				writer.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(writer).Encode(response)
			}
		} else {
			handler.ServeHTTP(writer, request)
		}
	}
}

func NotAuthenticated(writer http.ResponseWriter) {
	response := map[string]interface{}{
		"message": "Invalid authentication token",
		"data":    nil,
	}
	writer.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(writer).Encode(response)
}

func isAuthenticatedAdmin(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		authToken := request.Header.Get("Authorization")
		if authToken == "" {
			NotAuthenticated(writer)
		}

		// check token in redis
		_, isFound := db.GetAdminSession(authToken)
		if !isFound {
			tokenString := strings.Replace(authToken, "Bearer ", "", 1)
			adminID, err := utils.ParseAuthToken(tokenString)
			if err != nil {
				NotAuthenticated(writer)
			}

			var authAdmin models.Admin

			database := db.Database.DB
			err = database.Where("ID = ?", adminID).First(&authAdmin).Error
			if err != nil {
				NotAuthenticated(writer)
			}

			err = db.CreateAdminSession(authToken, authAdmin)
			if err == nil {
				handler.ServeHTTP(writer, request)
			} else {
				response := map[string]interface{}{
					"message": "An Error Occured",
					"data":    nil,
				}
				writer.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(writer).Encode(response)
			}
		} else {
			handler.ServeHTTP(writer, request)
		}
	}
}
