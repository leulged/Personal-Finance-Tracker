package main

import (
	"log"
	"os"

	"personal-finance-tracker/Delivery/router"
	"personal-finance-tracker/Infrastructure/db"
	"personal-finance-tracker/Infrastructure/mongodb"
	"personal-finance-tracker/Infrastructure/service"
	usecase "personal-finance-tracker/UseCase"
)

func main() {
	// Load environment variables
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // default for local dev
	}

	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "personal_finance_tracker" // Fixed: no spaces in database name
	}

	log.Printf("🔗 Connecting to MongoDB at: %s", mongoURI)
	log.Printf("📊 Database name: %s", dbName)

	// Connect to MongoDB
	client, err := db.ConnectToMongoDB(mongoURI)
	if err != nil {
		log.Fatalf("❌ MongoDB connection failed: %v", err)
	}

	// Test the connection
	err = client.Ping(nil, nil)
	if err != nil {
		log.Fatalf("❌ MongoDB ping failed: %v", err)
	}
	log.Println("✅ MongoDB connection successful!")

	// Close connection when app exits
	defer func() {
		if err := client.Disconnect(nil); err != nil {
			log.Printf("⚠️ Failed to close MongoDB connection: %v", err)
		}
	}()

	// Get database collection
	database := client.Database(dbName)
	userCollection := database.Collection("users")
	log.Printf("📁 Using collection: %s", userCollection.Name())

	// Initialize repositories
	userRepo := repository.NewUserRepository(userCollection)

	// Initialize services
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key"
	}
	jwtService := services.NewJWTService(jwtSecret)

	// Initialize use cases
	userUsecase := usecase.NewUserUsecase(userRepo)

	// Setup router with dependencies
	router := router.SetupRouter(userUsecase, jwtService)

	// Start server
	log.Println("🚀 Personal Finance Tracker API running on http://localhost:8080")
	log.Printf("📊 Connected to database: %s", dbName)
	
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}