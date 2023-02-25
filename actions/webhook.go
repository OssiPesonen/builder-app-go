package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WebhookAction(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
