package controllers

import (
	"context"
	"net/http"
	"time"

	"miniproject/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var appCollection *mongo.Collection

func InitAppCollection(db *mongo.Database) {
	appCollection = db.Collection("applications")
}

func SubmitApplication(c *gin.Context) {
	name := c.PostForm("name")
	role := c.PostForm("role")
	file, _ := c.FormFile("resume")

	uploadPath := "./backend/uploads/" + file.Filename
	c.SaveUploadedFile(file, uploadPath)

	app := models.Application{
		Name:       name,
		Role:       role,
		ResumePath: "/uploads/" + file.Filename,
		CreatedAt:  time.Now(),
	}

	_, err := appCollection.InsertOne(context.TODO(), app)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application submitted successfully"})
}

func ListApplications(c *gin.Context) {
	cursor, err := appCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve applications"})
		return
	}

	var apps []models.Application
	if err := cursor.All(context.TODO(), &apps); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode applications"})
		return
	}

	c.JSON(http.StatusOK, apps)
}

func GetApplicationByName(c *gin.Context) {
	name := c.Param("name")
	var app models.Application
	err := appCollection.FindOne(context.TODO(), bson.M{"name": name}).Decode(&app)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	c.JSON(http.StatusOK, app)
}
