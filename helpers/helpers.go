package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sunset-wallpaper-changer-go/objects"
	"time"
)

const timeLiteral = "3:04:05 PM"

func AddHoursToTimeString(timeStr string, hours int) (string, error) {
	// Parse the time string into a time.Time value
	t, err := time.Parse(timeLiteral, timeStr)
	if err != nil {
		return "", err
	}

	// Add the desired number of hours
	t = t.Add(time.Duration(hours) * time.Hour)

	// Format the resulting time value back into a string
	return t.Format(timeLiteral), nil
}

// ParseTimeTwentyFour Parses time from format "3:04:05 PM" to 24-hour format
func ParseTimeTwentyFour(timeStr string) (string, error) {
	// Define the input time format
	inLayout := timeLiteral
	// Parse the input time string
	t, err := time.Parse(inLayout, timeStr)
	if err != nil {
		return "", err
	}
	// Define the output time format
	outLayout := "15:04"
	// Convert the time to the output format
	return t.Format(outLayout), nil
}

func ClosestToNow(sunrise time.Time, sunset time.Time) time.Time {
	now := time.Now()

	if now.Before(sunset) && now.After(sunrise) {
		res := time.Date(now.Year(), now.Month(), now.Day(), sunset.Hour(), sunset.Minute(), sunset.Second(), 0, now.Location())
		return res
	} else {
		res := time.Date(now.Year(), now.Month(), now.Day(), sunrise.Hour(), sunrise.Minute(), sunrise.Second(), 0, now.Location())
		res = res.AddDate(0, 0, 1)
		return res
	}
}

func ParseStringToTime(layout string, timeStr string) time.Time {
	t, err := time.Parse(layout, timeStr)
	now := time.Now()
	nowTime := time.Date(0000, 01, 01, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location())
	if err != nil {
		panic(err)
	}
	if nowTime.Before(t) {
		t = t.AddDate(now.Year(), int(now.Month())-1, now.Day())
	} else {
		t = t.AddDate(now.Year(), int(now.Month())-1, now.Day())
		t = t.AddDate(0, 0, 1)
	}
	return t
}

// GetResponseFromSunsetApi gets response from api based on user defined latitude and longitude in type Response
func GetResponseFromSunsetApi(latStr string, longStr string, offsetHours int) objects.Response {
	apiUrl := fmt.Sprintf("https://api.sunrise-sunset.org/json?lat=%s&lng=%s", latStr, longStr)

	response, err := http.Get(apiUrl)
	if err != nil {
		fmt.Print(err.Error() + " - fetching api was not successful")
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var parsedResponse objects.Response
	err = json.Unmarshal([]byte(string(responseData)), &parsedResponse)
	if err != nil {
		panic(err)
	}

	parsedResponse.Results.Sunrise, err = AddHoursToTimeString(parsedResponse.Results.Sunrise, offsetHours)
	if err != nil {
		panic(err)
	}
	parsedResponse.Results.Sunset, err = AddHoursToTimeString(parsedResponse.Results.Sunset, offsetHours)
	if err != nil {
		panic(err)
	}

	return parsedResponse
}

func SetWallpaperManjaroKDE(dirPtr string, sunset bool) {
	wallpaperPath := dirPtr
	if sunset {
		wallpaperPath += "night.jpg"
	} else {
		wallpaperPath += "day.jpg"
	}
	cmd := exec.Command("qdbus", "org.kde.plasmashell", "/PlasmaShell", "org.kde.PlasmaShell.evaluateScript",
		fmt.Sprintf("var Desktops = desktops();for (i=0;i<Desktops.length;i++) {d = Desktops[i]; d.wallpaperPlugin = \"org.kde.image\"; d.currentConfigGroup = Array(\"Wallpaper\", \"org.kde.image\", \"General\"); d.writeConfig(\"Image\", \"%s\")}", wallpaperPath))
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Output:", string(out))
	}
}

func SetExecutingWallpaperChangeToSunsetAndSunrise(parsedTime string) bool {
	// Set the time to execute the command (in 24-hour format)
	// Use the "at" command to schedule the command for execution
	atCmd := fmt.Sprintf("at %s", parsedTime)
	at := exec.Command("sh", "-c", atCmd)

	// Run the "at" command
	if err := at.Run(); err != nil {
		fmt.Println("Error:", err)
		return false
	}
	fmt.Println("Your wallpaper will be changed at", parsedTime)
	return true
}

func Log(text string) error {
	filename := "log.txt"
	// Convert the string to a byte slice
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error should be handled")
		}
	}(file)

	// Write the string to the file
	_, err = file.WriteString(text + "\n")
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

func RemoveAllScheduledAtCommands() {
	cut := exec.Command("cut", "-f1")
	atq := exec.Command("atq")

	// Get atq's stdout and attach it to cut's stdin.
	pipe, _ := atq.StdoutPipe()
	defer func(pipe io.ReadCloser) {
		err := pipe.Close()
		if err != nil {
			log.Println("Closing pipe was not successful")
		}
	}(pipe)

	cut.Stdin = pipe

	// Run atq first.
	err := atq.Start()
	if err != nil {
		return
	}

	res, _ := cut.Output()
	split := strings.Split(strings.Trim(string(res), "\n"), "\n")

	for _, val := range split {
		removeAt := exec.Command("atrm", val)
		err = removeAt.Run()
		if err != nil {
			return
		}
	}
}
