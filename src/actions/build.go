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
		c.JSON(http.StatusNotImplemented, gin.H{"description": "Build command missing. Nothing to execute."})
		return
	}

	// Use defined path from env, otherwise static
	lockFile := os.Getenv("BUILDER_LOCKFILE_PATH")

	if lockFile == "" {
		lockFile = os.TempDir() + "/builderLockFile"
	}

	// Check lockfile
	if _, err := os.Stat(lockFile); err == nil {
		c.JSON(http.StatusConflict, gin.H{"description": "Lockfile found. Build in progress. Please wait a moment."})
		return
	}

	// Run the
	go services.RunScript(filePath, lockFile)

	c.Writer.WriteHeader(200)

	return
}
