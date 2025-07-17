package controllers

import (
	"context"
	"net/http"
	"time"

	"fmt"                // Import fmt for error messages
	"miniproject/models" // Ensure this import path is correct based on your module

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options" // Import options for FindOne
)

var appCollection *mongo.Collection

func InitAppCollection(db *mongo.Database) {
	appCollection = db.Collection("applications")
}

func SubmitApplication(c *gin.Context) {
	name := c.PostForm("name")
	role := c.PostForm("role")

	// Validate required fields
	if name == "" || role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and Role are required fields."})
		return
	}

	file, err := c.FormFile("resume")
	var resumePath string
	if err != nil {
		// If resume is optional, handle the error gracefully.
		// If resume is mandatory, return an error.
		if err == http.ErrMissingFile {
			// Resume file is missing, but maybe it's optional.
			// For now, we'll allow it to be missing, but you can change this.
			fmt.Println("Resume file not provided (optional).")
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to get resume file: %v", err)})
			return
		}
	} else {
		uploadPath := "./backend/uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save resume file: %v", err)})
			return
		}
		resumePath = "/uploads/" + file.Filename
	}

	app := models.Application{
		Name:       name,
		Role:       role,
		ResumePath: resumePath, // Use the determined resumePath
		CreatedAt:  time.Now(),
	}

	// Use a context with timeout for database operations
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = appCollection.InsertOne(ctx, app)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save application: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application submitted successfully"})
}

func ListApplications(c *gin.Context) {
	// Use a context with timeout for database operations
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := appCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve applications: %v", err)})
		return
	}
	defer cursor.Close(ctx) // Important: Close the cursor to release resources

	var apps []models.Application
	if err := cursor.All(ctx, &apps); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to decode applications: %v", err)})
		return
	}

	c.JSON(http.StatusOK, apps)
}

func GetApplicationByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application name is required."})
		return
	}

	// Use a context with timeout for database operations
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var app models.Application
	err := appCollection.FindOne(ctx, bson.M{"name": name}, options.FindOne()).Decode(&app) // Added options for consistency
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Application with name '%s' not found.", name)})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve application: %v", err)})
		}
		return
	}

	c.JSON(http.StatusOK, app)
}
