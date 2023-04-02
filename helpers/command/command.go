package command

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sunset-wallpaper-changer-go/logger"
	"time"
)

var LOGGER = logger.GetInstance().Logger

func SetExecutingWallpaperChangeToSunsetAndSunrise(timeVal time.Time) bool {
	// Set the time to execute the command (in 24-hour format)
	// Use the "at" command to schedule the command for execution
	parsedTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", timeVal.String())
	if err != nil {
		LOGGER.Println("Error parsing time")
		return false
	}
	atCmd := fmt.Sprintf("at %s", parsedTime.Format("15:04"))
	at := exec.Command("sh", "-c", atCmd)

	// Run the "at" command
	if err := at.Run(); err != nil {
		LOGGER.Printf("Error: %s \n", err)
		fmt.Println(at.String())
		return false
	}

	fmt.Printf("Your wallpaper will be changed at %s\n", timeVal.String())
	return true
}

func RemoveAllScheduledAtCommands() {
	var (
		cut  *exec.Cmd
		atq  *exec.Cmd
		pipe io.ReadCloser
		res  []byte
	)

	cut = exec.Command("cut", "-f1")
	atq = exec.Command("atq")

	// Get atq's stdout and attach it to cut's stdin.
	pipe, _ = atq.StdoutPipe()
	defer func(pipe io.ReadCloser) {
		err := pipe.Close()
		if err != nil {
			LOGGER.Println("Closing pipe was not successful")
		}
	}(pipe)

	cut.Stdin = pipe

	// Run atq first.
	err := atq.Start()
	if err != nil {
		LOGGER.Printf("executing command atq was not successful - %s", err.Error())
		return
	}

	res, _ = cut.Output()
	split := strings.Split(strings.Trim(string(res), "\n"), "\n")
	if len(res) != 0 {
		for _, val := range split {
			removeAt := exec.Command("atrm", val)
			err = removeAt.Run()
			if err != nil {
				LOGGER.Printf("executing command atrm was not successful - %s", err.Error())
				return
			}
		}
	}
}
