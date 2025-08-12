package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"personal-finance-tracker/Infrastructure/db"
	
)

func main() {
	// Load environment variables
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // default for local dev
	}

	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "Personal Finance Tracker"
	}

	// Connect to MongoDB
	client, err := db.ConnectToMongoDB(mongoURI)
	if err != nil {
		log.Fatalf("❌ MongoDB connection failed: %v", err)
	}

	// Close connection when app exits
	defer func() {
		if err := client.Disconnect(nil); err != nil {
			log.Printf("⚠️ Failed to close MongoDB connection: %v", err)
		}
	}()

	// Initialize Gin
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "API is running and connected to MongoDB",
		})
	})

	// Example route to list DB name
	router.GET("/db-info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"database": dbName,
		})
	})

	// Start server
	log.Println("🚀 Server running on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}