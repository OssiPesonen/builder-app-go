package main

import (
	"builder-app/actions"
	"builder-app/services"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	r := gin.Default()
	loadRoutes(r)
	r.Run()
}

// Load env variables from .env file
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// Defined routes
func loadRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		// This is just a dummy to see if our app is up
		c.JSON(http.StatusOK, gin.H{"hello": "world"})
	})

	r.POST("/"+os.Getenv("BUILDER_WEBHOOK_ROUTE"), services.AuthenticateWebhook(), actions.WebhookAction)
	r.POST("/"+os.Getenv("BUILDER_GITHUB_ROUTE"), services.AuthenticateGithubWebhook(), actions.GithubAction)
}
