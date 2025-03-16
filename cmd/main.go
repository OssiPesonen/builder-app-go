package main

import (
	"fmt"
	"log"
	"os"

	"github.com/OssiPesonen/builder-app-go/internal/actions"
	"github.com/OssiPesonen/builder-app-go/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	// Delete lockfile when initializing
	lockFile := os.Getenv("BUILDER_LOCKFILE_PATH")
	os.Remove(lockFile)

	r := setupRoutes()
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

func setupRoutes() *gin.Engine {
	r := gin.Default()

	r.GET("/healthcheck", func(c *gin.Context) {
		c.Writer.WriteHeader(200)
	})

	r.POST("/build", middleware.AuthenticateRequest(), actions.BuildAction)

	return r
}
