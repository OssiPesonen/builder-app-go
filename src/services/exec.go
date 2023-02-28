package services

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/OssiPesonen/builder-app-go/src/utils"
)

func RunScript(filePath, lockFile string) {
	file, e := os.Create(lockFile)

	if e != nil {
		fmt.Printf("Unable to create lockfile")
		fmt.Printf(e.Error())
	}

	file.Close()

	// Setup a log streamer for some improved logging
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	logStreamerOut := utils.NewLogstreamer(logger, "stdout", false)
	defer logStreamerOut.Close()

	logStreamerErr := utils.NewLogstreamer(logger, "stderr", true)
	defer logStreamerErr.Close()

	cmd := exec.Command(filePath)
	cmd.Stdout = logStreamerOut
	cmd.Stderr = logStreamerErr

	logStreamerErr.FlushRecord()

	err := cmd.Start()

	if err != nil {
		fmt.Printf(err.Error())
	}

	// And when you need to wait for the command to finish:
	if err := cmd.Wait(); err == nil {
		os.Remove(lockFile)
	}
}
