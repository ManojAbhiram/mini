package main

import (
	"context"
	"log"
	"miniproject/controllers"
	"time" // Import time for context timeout

	"github.com/gin-contrib/cors" // Recommended for robust CORS
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := gin.Default()

	// CORS middleware - Recommended to use gin-contrib/cors for better configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}                                       // Allow all origins for development, restrict in production
	config.AllowMethods = []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"} // Add other methods if needed
	config.AllowHeaders = []string{"Content-Type", "Authorization"}           // Add other headers if needed
	r.Use(cors.New(config))

	// Establish MongoDB connection with a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure the context is cancelled to release resources

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err) // Use Fatalf to exit if connection fails
	}

	// Ping the MongoDB server to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Successfully connected to MongoDB!")

	db := client.Database("job_portal")
	controllers.InitAppCollection(db)

	r.Static("/uploads", "./backend/uploads")

	r.POST("/api/apply", controllers.SubmitApplication)
	r.GET("/api/applications", controllers.ListApplications)
	r.GET("/api/application/:name", controllers.GetApplicationByName)

	// Graceful shutdown (optional, but good practice for production)
	// Example: Add a go routine to listen for OS signals and disconnect from MongoDB
	// This makes your server more resilient to shutdowns.

	log.Println("Server running on :8080")
	r.Run(":8080")
}
