package main

import (
	"log"
	"myapi/db"
	"myapi/routes"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize the database
	db.Init()

// Get the port from environment variables or use default port 8080
port := os.Getenv("PORT")
if port == "" {
	port = "8080" // Default port
}

// Initialize the router
router := routes.InitializeRoutes()

// Start the server
log.Printf("Starting server on port %s...\n", port)
if err := http.ListenAndServe(":"+port, router); err != nil {
	log.Fatalf("Server failed to start: %v", err)
}
}
