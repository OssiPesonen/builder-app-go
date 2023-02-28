package actions

import (
	"net/http"
	"os"

	"github.com/OssiPesonen/builder-app-go/src/services"
	"github.com/gin-gonic/gin"
)

func BuildAction(c *gin.Context) {
	filePath := os.Getenv("BUILDER_EXEC_PATH")

	// Check executable
	if filePath == "" {
		c.JSON(http.StatusNotImplemented, gin.H{"status": "error", "error": "Build command missing. Nothing to execute."})
		return
	}

	lockFile := os.TempDir() + "/builderLockFile"

	// Check lockfile
	if _, err := os.Stat(lockFile); err == nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "error": "Lockfile found. Build in progress. Please wait a moment."})
		return
	}

	go services.RunScript(filePath, lockFile)
	c.JSON(http.StatusAccepted, gin.H{"status": "ok"})
	return
}
