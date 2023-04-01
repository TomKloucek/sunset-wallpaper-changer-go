package kde

import (
	"strconv"
	"sunset-wallpaper-changer-go/helpers"
	"time"
)

const timeLiteral = " - Current time: "

func WallpaperChanger(latPtr float64, longPtr float64, wallpaperDir string) {
	// Getting current time and size of offset in location user provided
	currentTime := time.Now()
	_, offset := currentTime.Zone()
	offsetHours := offset / 3600

	// Truncate values because we do not need so precise values
	latStr := strconv.FormatFloat(latPtr, 'f', 6, 64)
	longStr := strconv.FormatFloat(longPtr, 'f', 6, 64)

	parsedResponse := helpers.GetResponseFromSunsetApi(latStr, longStr, offsetHours)

	layout := "15:04"
	twentyFourSunset, _ := helpers.ParseTimeTwentyFour(parsedResponse.Results.Sunset)
	twentyFourSunrise, _ := helpers.ParseTimeTwentyFour(parsedResponse.Results.Sunrise)
	sunset := helpers.ParseStringToTime(layout, twentyFourSunset)
	sunrise := helpers.ParseStringToTime(layout, twentyFourSunrise)

	helpers.RemoveAllScheduledAtCommands()

	closerTime := helpers.ClosestToNow(sunrise, sunset)
	err := helpers.Log("Closest time: " + closerTime.String() + timeLiteral + time.Now().String())

	if closerTime == sunset {
		helpers.SetExecutingWallpaperChangeToSunsetAndSunrise(twentyFourSunset, true)
	} else {
		helpers.SetExecutingWallpaperChangeToSunsetAndSunrise(twentyFourSunrise, true)
	}
	if currentTime.Before(sunset) && currentTime.After(sunrise) {
		helpers.SetWallpaperManjaroKDE(wallpaperDir, false)
		err := helpers.Log("Setting wallpaper to day.jpg" + timeLiteral + time.Now().String())
		if err != nil {
			return
		}
	} else {
		helpers.SetWallpaperManjaroKDE(wallpaperDir, true)
		err := helpers.Log("Setting wallpaper to night.jpg" + timeLiteral + time.Now().String())
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}
}
