package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/OssiPesonen/builder-app-go/src/actions"
	"github.com/OssiPesonen/builder-app-go/src/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	r := gin.Default()
	loadRoutes(r)
	r.Run(":" + os.Getenv("BUILDER_PORT"))
}

// Load env variables from .env file
func loadEnv() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println(err)
		log.Fatal("Error loading .env file")
	}
}

func loadRoutes(r *gin.Engine) {
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/webhook", services.AuthenticateWebhook(), actions.WebhookAction)
	r.POST("/github_webhook", services.AuthenticateGithubWebhook(), actions.GithubAction)
}
