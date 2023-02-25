package services

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Middleware to authenticate a custom webhook request
func AuthenticateWebhook() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := os.Getenv("BUILDER_WEBHOOK_SECRET")

		if secret == "" {
			c.Next()
		}

		headerSecret := c.Request.Header.Get(os.Getenv("BUILDER_WEBHOOK_SECRET_HEADER"))

		if headerSecret == secret {
			c.Next()
		} else {
			c.Abort()
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Invalid secret"})
		}
	}
}

// Middleware to authenticate Github commit event
func AuthenticateGithubWebhook() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := os.Getenv("BUILDER_GITHUB_SECRET")
		_, err := VerifyGithubReqSignature([]byte(secret), c.Request)

		if err == nil {
			c.Next()
		} else {
			log.Printf(err.Error())
			c.Abort()
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Unable to validate signature"})
		}
	}
}
