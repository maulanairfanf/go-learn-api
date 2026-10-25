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
	godotenv.Load()

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + database + " port=" + port + " sslmode=disable TimeZone=Asia/Jakarta"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(&models.Category{}, &models.Product{}, &models.User{}, &models.Order{}, &models.OrderItem{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}
}

func Close() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Error getting DB instance: %v", err)
		return
	}
	sqlDB.Close()
}
