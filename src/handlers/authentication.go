package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/codemunsta/risevest-test/src/db"
	"github.com/codemunsta/risevest-test/src/models"
	"github.com/codemunsta/risevest-test/src/utils"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type LoginStruct struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {

		var newUser models.User
		request.ParseForm()
		err := json.NewDecoder(request.Body).Decode(&newUser)
		if err != nil {
			http.Error(writer, "Invalid request", http.StatusBadRequest)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(writer, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		newUser.Password = string(hashedPassword)

		db.Database.DB.Create(&newUser)

		writer.WriteHeader(http.StatusCreated)
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(map[string]string{"message": "User registered successfully"})
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func RegisterAdmin(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		var newAdmin models.Admin

		request.ParseForm()
		err := json.NewDecoder(request.Body).Decode(&newAdmin)
		if err != nil {
			http.Error(writer, "Invalid request", http.StatusBadRequest)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newAdmin.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(writer, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		newAdmin.Password = string(hashedPassword)
		// newAdmin.Role = "admin"

		db.Database.DB.Create(&newAdmin)

		writer.WriteHeader(http.StatusCreated)
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(map[string]string{"message": "Admin registered successfully"})

	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func LoginAuthentication(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		var loginForm LoginStruct
		err := json.NewDecoder(request.Body).Decode(&loginForm)
		if err != nil {
			http.Error(writer, "Invalid request", http.StatusBadRequest)
			return
		}
		var user models.User

		db := db.Database.DB

		err = db.Where("email = ?", loginForm.Email).First(&user).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				http.Error(writer, "User not found", http.StatusNotFound)
				return
			}
			http.Error(writer, "Database error", http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password))
		if err != nil {
			http.Error(writer, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		authToken, err := utils.GenerateAuthToken(user.ID)
		if err != nil {
			http.Error(writer, "Failed to generate token", http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(map[string]string{"token": authToken})

	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func LoginAdmin(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		var loginForm LoginStruct
		err := json.NewDecoder(request.Body).Decode(&loginForm)
		if err != nil {
			http.Error(writer, "Invalid request", http.StatusBadRequest)
			return
		}
		var user models.Admin

		db := db.Database.DB

		err = db.Where("email = ?", loginForm.Email).First(&user).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				http.Error(writer, "User not found", http.StatusNotFound)
				return
			}
			http.Error(writer, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password))
		if err != nil {
			http.Error(writer, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		authToken, err := utils.GenerateAuthToken(user.ID)
		if err != nil {
			http.Error(writer, "Failed to generate token", http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(map[string]string{"token": authToken})
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
