package actions

import (
	"net/http"
	"os"

	"github.com/OssiPesonen/builder-app-go/src/services"
	"github.com/gin-gonic/gin"
)

func GithubAction(c *gin.Context) {
	services.RunScript(os.Getenv("BUILDER_EXEC_PATH"))
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
