package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/codemunsta/risevest-test/src/db"
	"github.com/codemunsta/risevest-test/src/models"
)

func AdminGetAllFiles(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		var files []models.File
		database := db.Database.DB
		err := database.Find(&files).Error
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

		fileIDS := request.FormValue("filesID")
		fileID, _ := strconv.Atoi(fileIDS)

		var file models.FileDelete
		database := db.Database.DB
		err := database.Where("file_id = ?", fileID).First(&file).Error
		if err != nil {
			file.FileID = uint(fileID)
			file.AdminApproval1 = authAdmin.ID
			database.Create(&file)
			response := map[string]interface{}{
				"message": "File marked as unsafe",
			}
			writer.WriteHeader(http.StatusOK)
			json.NewEncoder(writer).Encode(response)
		} else {
			if file.AdminApproval2 == 0 && file.AdminApproval1 != authAdmin.ID {
				file.AdminApproval2 = authAdmin.ID
				database.Save(&file)
				response := map[string]interface{}{
					"message": "File marked as unsafe",
				}
				writer.WriteHeader(http.StatusOK)
				json.NewEncoder(writer).Encode(response)
			} else if file.AdminApproval3 == 0 && file.AdminApproval2 != authAdmin.ID && file.AdminApproval1 != authAdmin.ID {
				file.AdminApproval3 = authAdmin.ID
				database.Save(&file)
				err := deleteFile(uint(fileID))
				if err != nil {
					response := map[string]interface{}{
						"message": "Marked but not deleted: " + err.Error(),
					}
					writer.WriteHeader(http.StatusOK)
					json.NewEncoder(writer).Encode(response)
				}
				response := map[string]interface{}{
					"message": "File marked as unsafe",
				}
				writer.WriteHeader(http.StatusOK)
				json.NewEncoder(writer).Encode(response)
			} else {
				response := map[string]interface{}{
					"message": "Failed to mark file",
				}
				writer.WriteHeader(http.StatusOK)
				json.NewEncoder(writer).Encode(response)
			}
		}
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func deleteFile(fileID uint) error {
	var file models.FileDelete
	database := db.Database.DB
	err := database.Where("file_id = ?", fileID).First(&file).Error
	if err != nil {
		return err
	} else {
		err := database.Delete(&file).Error
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}
