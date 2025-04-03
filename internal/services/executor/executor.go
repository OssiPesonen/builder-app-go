package executor

import (
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/OssiPesonen/builder-app-go/internal/services/broadcaster"
	"github.com/OssiPesonen/builder-app-go/internal/services/broadcaster/channels"
	"github.com/OssiPesonen/builder-app-go/internal/utils"
)

func RunScript(filePath, lockFile string) {
	channels, wg := createBroadcaster()

	// Setup a log streamer for some improved logging
	logger := log.New(os.Stdout, "[BUILDER] ", log.Ldate|log.Ltime)

	logStreamerOut := utils.NewLogstreamer(logger, "stdout", false)
	defer logStreamerOut.Close()

	logStreamerErr := utils.NewLogstreamer(logger, "stderr", true)
	defer logStreamerErr.Close()

	logStreamerOut.Logger.Println("Creating lock")
	file, e := os.Create(lockFile)

	if e != nil {
		logStreamerErr.Logger.Println("Unable to create lockfile")
		logStreamerErr.Logger.Println(e.Error())
	}

	file.Close()

	channels.Publish(
		broadcaster.Message{
			Title: "Initiating build",
			Body:  ":rocket: *Initiating build*\nReceived a request to start a new build. Website will be updated shortly.",
		},
	)

	cmd := exec.Command("bash", filePath)
	cmd.Stdout = logStreamerOut
	cmd.Stderr = logStreamerErr

	logStreamerErr.Logger.Println("Executing build command...")
	err := cmd.Start()
	cmd.Wait()

	if err != nil {
		logStreamerErr.Logger.Println(err.Error())
	}

	logStreamerErr.FlushRecord()
	// Finally remove lockfile
	os.Remove(lockFile)

	channels.Publish(
		broadcaster.Message{
			Title: "All done. Website should now be updated.",
			Body:  ":tada: *All done!*\nWebsite should now be updated. Enjoy your fresh new content!",
		},
	)

	wg.Wait()
	channels.Close()
}

func createBroadcaster() (*broadcaster.Broadcaster, *sync.WaitGroup) {
	// Create a message caster to publish messages to external channels
	caster := broadcaster.NewBroadcaster()
	enabledChannels := make(map[string]broadcaster.Subscriber)

	slackWebhookUrl := os.Getenv("SLACK_WEBHOOK_URL")
	if slackWebhookUrl != "" {
		slackChan := caster.Subscribe()
		enabledChannels["slack"] = channels.NewSlack(os.Getenv("SLACK_WEBHOOK_URL"), slackChan)
	}

	var wg sync.WaitGroup
	wg.Add(len(enabledChannels))

	for _, v := range enabledChannels {
		go v.Run(&wg)
	}

	return caster, &wg
}
