package main

import (
	"log"
	"os"

	"cafe-api/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load the .env file so os.Getenv() works
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Connect to the Database
	database.Connect()

	// 3. Initialize the Gin router
	router := gin.Default()

	// 4. Create a quick test route to make sure Gin is working
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong! The Café API is running.",
		})
	})

	// 5. Start the server
	port := os.Getenv("PORT")
	log.Printf("🚀 Server is running on port %s", port)
	router.Run(":" + port)
}
