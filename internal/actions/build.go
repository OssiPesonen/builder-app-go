package actions

import (
	"net/http"
	"os"
	"sync"

	"github.com/OssiPesonen/builder-app-go/internal/channels"
	"github.com/OssiPesonen/builder-app-go/internal/services"
	"github.com/gin-gonic/gin"
)

func BuildAction(c *gin.Context) {
	broadcaster, waitGroup := createBroadcaster()
	filePath := os.Getenv("BUILDER_EXEC_PATH")

	broadcaster.Publish(
		services.Message{
			Title: "Initiating build",
			Body:  ":rocket: *Initiating build*\nReceived a request to start a new build. Website will be updated shortly.",
		},
	)

	// Check executable
	if filePath == "" {
		broadcaster.Publish(
			services.Message{
				Title: "Build command missing",
				Body:  ":warning: *Build command missing*\nNothing to execute. Please add the build command to the environment variables on the server.",
			},
		)

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
		broadcaster.Publish(
			services.Message{
				Title: "Lockfile found. Build in progress. Please wait a moment.",
				Body:  ":lock: *Lockfile found*\nBuilding is currently in progress. Please wait a moment.",
			},
		)

		c.JSON(http.StatusConflict, gin.H{"description": "Lockfile found. Build in progress. Please wait a moment."})
		return
	}

	// Run the
	go services.RunScript(filePath, lockFile, broadcaster)

	waitGroup.Wait()
	c.Writer.WriteHeader(200)
}

func createBroadcaster() (*services.MessageBroadcaster, *sync.WaitGroup) {
	// Create a message broadcaster to publish messages to external channels
	broadcaster := services.NewMessageBroadcaster()
	enabledChannels := make(map[string]services.Subscriber)

	slackWebhookUrl := os.Getenv("SLACK_WEBHOOK_URL")
	if slackWebhookUrl != "" {
		slackChan := broadcaster.Subscribe()
		enabledChannels["slack"] = channels.NewSlack(os.Getenv("SLACK_WEBHOOK_URL"), slackChan)
	}

	var wg sync.WaitGroup
	wg.Add(len(enabledChannels))

	for _, v := range enabledChannels {
		go v.Run(&wg)
	}

	return broadcaster, &wg
}
