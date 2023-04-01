package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sunset-wallpaper-changer-go/kde"
	"sunset-wallpaper-changer-go/logger"
)

const dayFile string = "day.jpg"
const nightFile string = "night.jpg"

func main() {
	var (
		dirPtr  *string
		latPtr  *float64
		longPtr *float64
		LOGGER  *log.Logger
	)
	LOGGER = logger.GetInstance().Logger

	dirPtr = flag.String("dir", "/usr/share/backgrounds/", "directory from where we get day and night wallpapers")
	latPtr = flag.Float64("lat", 50.08032461927618, "Your current latitude")
	longPtr = flag.Float64("long", 14.430143915469639, "Your current longitude")

	flag.Parse()

	if _, err := os.Stat(*dirPtr + "/" + dayFile); err == nil {
		if _, err := os.Stat(*dirPtr + "/" + nightFile); err != nil {
			LOGGER.Printf("Mandatory file %s does not exist\n", nightFile)
			os.Exit(1)
		}
	} else {
		LOGGER.Printf("Mandatory file %s does not exist\n", dayFile)
		os.Exit(1)
	}

	switch runningOs := runtime.GOOS; runningOs {
	case "darwin":
		LOGGER.Println("User tried to use MacOS")
		fmt.Println("Program is unfortunately currently not working on MacOS")
	case "linux":
		kde.WallpaperChanger(*latPtr, *longPtr, *dirPtr)
	default:
		// freebsd, openbsd, plan9, windows
		LOGGER.Printf("User tried to use %s\n", runningOs)
		fmt.Printf("Program is unfortunately currently not working on %s\n", runningOs)
	}

}
