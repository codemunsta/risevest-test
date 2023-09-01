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
		"host=%s user=%s password=%s dbname=%s port=6247 sslmode=disable TimeZone=Africa/Lagos",
		os.Getenv("PGHOST"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGDATABASE"),
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
