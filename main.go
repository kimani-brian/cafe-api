package main

import (
	"log"
	"os"

	"cafe-api/database"
	"cafe-api/routes" // <-- Import our new routes package

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Connect()
	database.RunMigrations()

	// Use the router we just created!
	router := routes.SetupRouter()

	port := os.Getenv("PORT")
	log.Printf("🚀 Server is running on port %s", port)
	router.Run(":" + port)
}
