package main

import (
	"context"
	"log"
	"miniproject/controllers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := gin.Default()

	// CORS middleware (optional for dev)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("job_portal")
	controllers.InitAppCollection(db)

	r.Static("/uploads", "./backend/uploads")

	r.POST("/api/apply", controllers.SubmitApplication)
	r.GET("/api/applications", controllers.ListApplications)
	r.GET("/api/application/:name", controllers.GetApplicationByName)

	r.Run(":8080")
}
