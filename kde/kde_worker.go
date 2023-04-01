package kde

import (
	"fmt"
	"os/exec"
	"strconv"
	apiHelper "sunset-wallpaper-changer-go/helpers/api"
	commandHelper "sunset-wallpaper-changer-go/helpers/command"
	timeHelper "sunset-wallpaper-changer-go/helpers/time"
	"sunset-wallpaper-changer-go/logger"
	"time"
)

const dayFile string = "day.jpg"
const nightFile string = "night.jpg"

var LOGGER = logger.GetInstance().Logger

func setWallpaperManjaroKDE(dirPtr string, sunset bool) {
	var (
		wallpaperPath string
		cmd           *exec.Cmd
		out           []byte
		err           error
	)

	LOGGER = logger.GetInstance().Logger

	wallpaperPath = dirPtr
	if sunset {
		wallpaperPath += dayFile
	} else {
		wallpaperPath += nightFile
	}

	cmd = exec.Command("qdbus", "org.kde.plasmashell", "/PlasmaShell", "org.kde.PlasmaShell.evaluateScript",
		fmt.Sprintf("var Desktops = desktops();for (i=0;i<Desktops.length;i++) {d = Desktops[i]; d.wallpaperPlugin = \"org.kde.image\"; d.currentConfigGroup = Array(\"Wallpaper\", \"org.kde.image\", \"General\"); d.writeConfig(\"Image\", \"%s\")}", wallpaperPath))
	out, err = cmd.CombinedOutput()
	if err != nil {
		LOGGER.Println("Error:", err)
		LOGGER.Println("Output:", string(out))
	}
}

func WallpaperChanger(latPtr float64, longPtr float64, wallpaperDir string) {
	var (
		currentTime       time.Time
		offset            int
		offsetHours       int
		latStr            string
		longStr           string
		layout            string
		twentyFourSunset  string
		twentyFourSunrise string
		sunset            time.Time
		sunrise           time.Time
		err               error
		closerTime        time.Time
	)
	// Getting current time and size of offset in location user provided
	currentTime = time.Now()
	_, offset = currentTime.Zone()
	offsetHours = offset / 3600

	LOGGER = logger.GetInstance().Logger

	// Truncate values because we do not need so precise values
	latStr = strconv.FormatFloat(latPtr, 'f', 6, 64)
	longStr = strconv.FormatFloat(longPtr, 'f', 6, 64)

	parsedResponse := apiHelper.GetResponseFromSunsetApi(latStr, longStr, offsetHours)

	layout = "15:04"
	twentyFourSunset, err = timeHelper.ParseTimeTwentyFour(parsedResponse.Results.Sunset)
	if err != nil {
		LOGGER.Println(err.Error())
	}

	twentyFourSunrise, err = timeHelper.ParseTimeTwentyFour(parsedResponse.Results.Sunrise)
	if err != nil {
		LOGGER.Println(err.Error())
	}

	sunset = timeHelper.ParseStringToTime(layout, twentyFourSunset)
	sunrise = timeHelper.ParseStringToTime(layout, twentyFourSunrise)

	commandHelper.RemoveAllScheduledAtCommands()

	closerTime = timeHelper.ClosestToNow(sunrise, sunset)

	if closerTime == sunset {
		commandHelper.SetExecutingWallpaperChangeToSunsetAndSunrise(twentyFourSunset)
	} else {
		commandHelper.SetExecutingWallpaperChangeToSunsetAndSunrise(twentyFourSunrise)
	}
	if currentTime.Before(sunset) && currentTime.After(sunrise) {
		setWallpaperManjaroKDE(wallpaperDir, false)
		LOGGER.Printf("Setting wallpaper to %s %s", dayFile, time.Now().String())
	} else {
		setWallpaperManjaroKDE(wallpaperDir, true)
		LOGGER.Printf("Setting wallpaper to %s %s", nightFile, time.Now().String())
	}
}
