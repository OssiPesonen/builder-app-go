package middleware

import (
	"net/http"
	"os"

	"github.com/OssiPesonen/builder-app-go/internal/services"
	"github.com/gin-gonic/gin"
)

// Middleware to authenticate requests
func AuthenticateRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request

		if req.Header.Get("x-github-event") != "" {
			AuthenticateGithubWebhook(c)
		} else {
			AuthenticateWebhook(c)
		}
	}
}

// Authenticate a custom build request
func AuthenticateWebhook(c *gin.Context) {
	secret := os.Getenv("BUILDER_WEBHOOK_SECRET")

	if secret == "" {
		// We can skip authentication if no secret is provided by environment
		c.Next()
	} else {
		headerSecret := c.Request.Header.Get(os.Getenv("BUILDER_WEBHOOK_SECRET_HEADER"))

		if headerSecret == secret {
			c.Next()
		} else {
			c.Abort()
			// Secret does not match
			c.JSON(http.StatusForbidden, gin.H{"description": "Authentication failed."})
		}
	}
}

// Authenticate Github event
func AuthenticateGithubWebhook(c *gin.Context) {
	secret := os.Getenv("BUILDER_GITHUB_SECRET")
	_, err := services.VerifyGithubReqSignature([]byte(secret), c.Request)

	if err == nil {
		c.Next()
	} else {
		c.Abort()
		// Invalid secret
		c.JSON(http.StatusForbidden, gin.H{"description": "Authentication failed."})
	}
}
