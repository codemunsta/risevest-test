package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/codemunsta/risevest-test/src/db"
	"github.com/codemunsta/risevest-test/src/models"
	"github.com/codemunsta/risevest-test/src/utils"
)

type UploadStruct struct {
	FolderName string
}

func NewApi(writer http.ResponseWriter, request *http.Request) {
	response := map[string]interface{}{
		"User":    request.RemoteAddr,
		"Message": "Pong",
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}

func UploadFile(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {

		// Get authorization token
		authToken := request.Header.Get("Authorization")
		if authToken == "" {
			http.Error(writer, "Missing authentication token", http.StatusUnauthorized)
			return
		}

		// Get request body
		var uploadBody UploadStruct
		request.ParseForm()
		err := json.NewDecoder(request.Body).Decode(&uploadBody)
		if err != nil {
			http.Error(writer, "Invalid request", http.StatusBadRequest)
			return
		}

		// Verify authorization token
		tokenString := strings.Replace(authToken, "Bearer ", "", 1)
		userID, err := utils.ParseAuthToken(tokenString)
		if err != nil {
			http.Error(writer, "Invalid authentication token", http.StatusUnauthorized)
			return
		}

		// fetch authorized user
		var authuser models.User

		db := db.Database.DB
		err = db.Where("ID = ?", userID).First(&authuser).Error
		if err != nil {
			http.Error(writer, "Invalid authentication token", http.StatusUnauthorized)
			return
		}

		// fetch file
		file, fileHeader, err := request.FormFile("file")
		if err != nil {
			http.Error(writer, "Could not process file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// return files greater than 200mb
		if fileHeader.Size > 200*1024*1024 {
			http.Error(writer, "File size exceeds 200 MB", http.StatusBadRequest)
			return
		}
		filePath := authuser.Email + "uploads/" + uploadBody.FolderName + fileHeader.Filename

		// create file model
		outFile, err := os.Create(filePath)
		if err != nil {
			http.Error(writer, "Failed to create file", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		if err != nil {
			http.Error(writer, "Failed to copy file contents", http.StatusInternalServerError)
			return
		}

		var uFile models.File
		uFile.FolderName = uploadBody.FolderName
		uFile.FileName = fileHeader.Filename
		uFile.FilePath = filePath
		uFile.UserID = userID

		db.Create(&uFile)

		writer.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"Message": "File uploaded successfully",
		}
		json.NewEncoder(writer).Encode(response)
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func CreateFolder(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		// Get authorization token
		authToken := request.Header.Get("Authorization")
		if authToken == "" {
			http.Error(writer, "Missing authentication token", http.StatusUnauthorized)
			return
		}

		// Verify authorization token
		tokenString := strings.Replace(authToken, "Bearer ", "", 1)
		userID, err := utils.ParseAuthToken(tokenString)
		if err != nil {
			http.Error(writer, "Invalid authentication token", http.StatusUnauthorized)
			return
		}

		// fetch authorized user
		var authuser models.User

		db := db.Database.DB
		err = db.Where("ID = ?", userID).First(&authuser).Error
		if err != nil {
			http.Error(writer, "Invalid authentication token", http.StatusUnauthorized)
			return
		}

		// Get request body
		var folderBody models.Folder
		folderBody.UserID = userID
		folderBody.User = authuser
		request.ParseForm()
		err = json.NewDecoder(request.Body).Decode(&folderBody)
		if err != nil {
			http.Error(writer, "Invalid request", http.StatusBadRequest)
			return
		}

		db.Create(&folderBody)
		writer.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"Message": "Folder created successfully",
		}
		json.NewEncoder(writer).Encode(response)
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
