package actions

import (
	"net/http"
	"os"

	"github.com/OssiPesonen/builder-app-go/src/services"
	"github.com/gin-gonic/gin"
)

func BuildAction(c *gin.Context) {
	filePath := os.Getenv("BUILDER_EXEC_PATH")

	if filePath == "" {
		c.JSON(http.StatusOK, gin.H{"status": "error", "error": "Build command missing. Nothing to execute."})
	} else {
		services.RunScript(filePath)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
