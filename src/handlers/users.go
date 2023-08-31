package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/codemunsta/risevest-test/src/db"
	"github.com/codemunsta/risevest-test/src/models"
)

func UploadFile(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {

		// Get authorized User
		authToken := request.Header.Get("Authorization")
		authenticatedUser, _ := db.GetSession(authToken)
		authUser := authenticatedUser.User

		// Get request body
		err := request.ParseMultipartForm(200 * 1024 * 1024)
		if err != nil {
			http.Error(writer, "Error parsing form", http.StatusBadRequest)
			return
		}

		// fetch file
		file, fileHeader, err := request.FormFile("file")
		if err != nil {
			http.Error(writer, "Could not process file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		name := request.FormValue("name")

		var folder models.Folder
		database := db.Database.DB
		err = database.Where("name = ? ANd user_id = ?", name, authUser.ID).First(&folder).Error
		if err != nil {
			http.Error(writer, "Folder not created", http.StatusNotFound)
			return
		}

		filePath := authUser.Email + "/uploads/" + folder.Name
		// Extract directory path from filePath
		dirPath := filepath.Dir(filePath)

		// Create directories if they don't exist
		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			http.Error(writer, "Failed to create directories: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Create the file
		outFile, err := os.Create(filePath)
		if err != nil {
			http.Error(writer, "Failed to create file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer os.Remove(outFile.Name())
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		if err != nil {
			http.Error(writer, "Error copying file content", http.StatusInternalServerError)
			return
		}

		// // upload to cloudinary

		var uFile models.File
		uFile.FolderName = folder.Name
		uFile.FileName = fileHeader.Filename
		uFile.UserID = authUser.ID
		uFile.Safe = true

		database.Create(&uFile)

		completedRoutine := make(chan bool)
		go uploadToCloudinary(outFile.Name(), uFile.ID, completedRoutine)

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
		// Get authorized User
		authToken := request.Header.Get("Authorization")
		authenticatedUser, _ := db.GetSession(authToken)
		authUser := authenticatedUser.User

		// Get request body
		var folderBody models.Folder
		request.ParseForm()
		err := json.NewDecoder(request.Body).Decode(&folderBody)
		if err != nil {
			http.Error(writer, "Invalid request", http.StatusBadRequest)
			return
		}
		folderBody.UserID = authUser.ID

		db.Database.DB.Create(&folderBody)
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

func DownloadFileCloudinary(writer http.ResponseWriter, request *http.Request) {
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
		err = db.Where("id = ?", fileID).First(&file).Error
		if err != nil {
			http.Error(writer, "File not found", http.StatusNotFound)
		}

		// download
		http.Redirect(writer, request, file.FilePath, http.StatusFound)
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func uploadToCloudinary(filePath string, id uint, completed chan bool) {
	cld, _ := cloudinary.NewFromURL(fmt.Sprintf("cloudinary://%v:%v@%v", os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_SECRET_KEY"), os.Getenv("CLOUDINARY_HOST")))
	var ctx = context.Background()
	cloudinaryUploader := cld.Upload
	response, err := cloudinaryUploader.Upload(ctx, filePath, uploader.UploadParams{})
	fmt.Print(response)
	if err != nil {
		log.Println("Error uploading file to Cloudinary:", err)
		return
	} else {
		var file models.File
		database := db.Database.DB
		err = database.Where("id = ?", id).First(&file).Error
		if err != nil {
			log.Fatal("an error updating file path occured")
			return
		} else {
			file.FilePath = response.SecureURL
			database.Save(file)
		}
	}
	fmt.Println("File uploaded to Cloudinary successfully")
	completed <- true
}

// get folders

func GetFolders(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {

		// Get authorized User
		authToken := request.Header.Get("Authorization")
		authenticatedUser, _ := db.GetSession(authToken)
		authUser := authenticatedUser.User

		// Get folders
		var folders []models.Folder
		database := db.Database.DB
		err := database.Where("user_id = ?", authUser.ID).Find(&folders).Error
		if err != nil {
			http.Error(writer, "Folder not created", http.StatusNotFound)
			return
		}
		response := map[string]interface{}{
			"message": "Here are your folders",
			"data":    folders,
		}
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(response)
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// get files

func GetFiles(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {

		// Get authorized User
		authToken := request.Header.Get("Authorization")
		authenticatedUser, _ := db.GetSession(authToken)
		authUser := authenticatedUser.User

		// Get folders
		var files []models.File
		database := db.Database.DB
		err := database.Where("user_id = ?", authUser.ID).Find(&files).Error
		if err != nil {
			http.Error(writer, "An error occured", http.StatusInternalServerError)
			return
		}
		response := map[string]interface{}{
			"message": "Here are your files",
			"data":    files,
		}
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(response)
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
