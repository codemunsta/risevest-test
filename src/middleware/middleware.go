package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/codemunsta/risevest-test/src/db"
	"github.com/codemunsta/risevest-test/src/models"
	"github.com/codemunsta/risevest-test/src/utils"
)

func IsAuthenticated(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		authToken := request.Header.Get("Authorization")
		if authToken == "" {
			NotAuthenticated(writer)
		}

		// check token in redis
		_, isFound := db.GetSession(authToken)
		if !isFound {
			log.Println("Hi didn't find redis session")
			tokenString := strings.Replace(authToken, "Bearer ", "", 1)
			userID, err := utils.ParseAuthToken(tokenString)
			if err != nil {
				NotAuthenticated(writer)
			}

			log.Println("Hi user authenticated")
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
				log.Println(err)
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

func IsAuthenticatedAdmin(handler http.HandlerFunc) http.HandlerFunc {
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
