package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codemunsta/risevest-test/src/db"
	"github.com/codemunsta/risevest-test/src/models"
)

func AdminGetAllFiles(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		var files []models.File
		database := db.Database.DB
		err := database.Where("user_id = ?").Find(&files).Error
		if err != nil {
			http.Error(writer, "No file available", http.StatusNotFound)
			return
		}
		response := map[string]interface{}{
			"message": "Here are your folders",
			"data":    files,
		}
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(response)
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func AdminMarkFileUnsafe(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {

		authToken := request.Header.Get("Authorization")
		authenticatedAdmin, _ := db.GetAdminSession(authToken)
		authAdmin := authenticatedAdmin.Admin

		fmt.Print(authAdmin)
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// func DeleteFile(file models.File) error {

// }
