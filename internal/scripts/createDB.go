package main

import (
	"fmt"
	"log"
	"os"

	"github.com/XineAurora/instantnotes-server/internal/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USERNAME"),
		os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB_NAME"),
		os.Getenv("POSTGRES_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Group{})
	db.AutoMigrate(&models.GroupMember{})
	db.AutoMigrate(&models.Note{})
	db.AutoMigrate(&models.Folder{})
	db.AutoMigrate(&models.FolderLink{})
}
