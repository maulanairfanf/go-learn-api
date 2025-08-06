package db

import (
	"log"
	"myapi/models"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

   // Get database credentials from environment variables
   user := os.Getenv("DB_USER")
   password := os.Getenv("DB_PASSWORD")
   host := os.Getenv("DB_HOST")
   port := os.Getenv("DB_PORT")
   database := os.Getenv("DB_NAME")

   dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + database + " port=" + port + " sslmode=disable TimeZone=Asia/Jakarta"
   DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
   if err != nil {
	   log.Fatal("Failed to connect to database:", err)
   }

	// Migrate the schema
	err = DB.AutoMigrate(&models.Category{}, &models.Product{}, &models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}
}
