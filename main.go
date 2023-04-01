package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sunset-wallpaper-changer-go/kde"
)

func main() {
	dirPtr := flag.String("dir", "/usr/share/backgrounds/", "directory from where we get day and night wallpapers")
	latPtr := flag.Float64("lat", 50.08032461927618, "Your current latitude")
	longPtr := flag.Float64("long", 14.430143915469639, "Your current longitude")

	flag.Parse()

	if _, err := os.Stat(*dirPtr + "/day.jpg"); err == nil {
		if _, err := os.Stat(*dirPtr + "/night.jpg"); err != nil {
			fmt.Printf("Mandatory file night.jpg does not exist\n")
			os.Exit(1)
		}
	} else {
		fmt.Printf("Mandatory file day.jpg does not exist\n")
		os.Exit(1)
	}

	switch runningOs := runtime.GOOS; runningOs {
	case "darwin":
		fmt.Println("Program is unfortunately currently not working on MacOS")
	case "linux":
		kde.WallpaperChanger(*latPtr, *longPtr, *dirPtr)
	default:
		// freebsd, openbsd, plan9, windows
		fmt.Printf("Program is unfortunately currently not working on %s\n", runningOs)
	}

}
