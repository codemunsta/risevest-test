package db

import (
	"fmt"
	"log"
	"os"

	"github.com/codemunsta/risevest-test/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBInstance struct {
	DB *gorm.DB
}

var Database DBInstance

func InitDB() {
	fmt.Println("Hey there")
	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Africa/Lagos",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Postgres failed to setup or connect. \n", err)
	}

	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migrations")
	db.AutoMigrate(&models.User{}, &models.Admin{}, &models.Folder{}, &models.File{}, &models.FileDelete{})

	Database = DBInstance{
		DB: db,
	}
}
