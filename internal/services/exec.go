package services

import (
	"log"
	"os"
	"os/exec"

	"github.com/OssiPesonen/builder-app-go/internal/utils"
)

func RunScript(filePath, lockFile string, messageBroadcaster *MessageBroadcaster) {
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

	logger.Println("Executing build...")

	cmd := exec.Command(filePath)
	cmd.Stdout = logStreamerOut
	cmd.Stderr = logStreamerErr

	logStreamerErr.FlushRecord()

	err := cmd.Start()

	if err != nil {
		logStreamerErr.Logger.Println(err.Error())
	} else {
		logStreamerOut.Logger.Println("Executing build...")
	}

	// Finally remove lockfile
	os.Remove(lockFile)
	messageBroadcaster.Publish(
		Message{
			Title: "All done. Website should now be updated.",
			Body:  ":tada: *All done!*\nWebsite should now be updated. Enjoy your fresh new content!",
		},
	)

	messageBroadcaster.Close()
}
