package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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

func FileDownload(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		queryParams := request.URL.Query()
		fileIDS := queryParams.Get("fileID")
		fileID, err := strconv.Atoi(fileIDS)
		if err != nil {
			http.Error(writer, "Invalid file ID", http.StatusBadRequest)
			return
		}
		// Fetch file information from the database
		var file models.File
		db := db.Database.DB
		if err := db.First(&file, fileID).Error; err != nil {
			http.Error(writer, "File not found", http.StatusNotFound)
			return
		}
		// Open the file on the server
		fileReader, err := os.Open(file.FilePath)
		if err != nil {
			http.Error(writer, "Failed to open file", http.StatusInternalServerError)
			return
		}
		defer fileReader.Close()

		// Set appropriate headers for file download

		// Get file info
		fileInfo, err := fileReader.Stat()
		if err != nil {
			http.Error(writer, "Failed to get file info", http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.FileName))
		writer.Header().Set("Content-Type", "application/octet-stream")
		writer.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

		// Stream the file's contents to the client's response
		_, err = io.Copy(writer, fileReader)
		if err != nil {
			http.Error(writer, "Failed to send file", http.StatusInternalServerError)
			return
		}
	}
}
