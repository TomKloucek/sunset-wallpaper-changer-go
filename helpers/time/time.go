package time

import (
	"sunset-wallpaper-changer-go/logger"
	"time"
)

const timeLiteral = "3:04:05 PM"

var LOGGER = logger.GetInstance().Logger

func AddHoursToTimeString(timeStr string, hours int) (string, error) {
	var (
		t   time.Time
		err error
	)
	// Parse the time string into a time.Time value
	t, err = time.Parse(timeLiteral, timeStr)
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
	var (
		inLayout  string
		t         time.Time
		outLayout string
		err       error
	)

	// Define the input time format
	inLayout = timeLiteral
	// Parse the input time string
	t, err = time.Parse(inLayout, timeStr)
	if err != nil {
		return "", err
	}
	// Define the output time format
	outLayout = "15:04"
	// Convert the time to the output format
	return t.Format(outLayout), nil
}

func ClosestToNow(sunrise time.Time, sunset time.Time) time.Time {
	var (
		now time.Time
		res time.Time
	)

	now = time.Now()

	if now.Before(sunset) && now.After(sunrise) {
		res = time.Date(now.Year(), now.Month(), now.Day(), sunset.Hour(), sunset.Minute(), sunset.Second(), 0, now.Location())
		return res
	} else {
		res = time.Date(now.Year(), now.Month(), now.Day(), sunrise.Hour(), sunrise.Minute(), sunrise.Second(), 0, now.Location())
		res = res.AddDate(0, 0, 1)
		return res
	}
}

func ParseStringToTime(layout string, timeStr string) time.Time {
	var (
		t       time.Time
		now     time.Time
		nowTime time.Time
		err     error
	)

	t, err = time.Parse(layout, timeStr)
	now = time.Now()
	nowTime = time.Date(0000, 01, 01, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location())
	if err != nil {
		LOGGER.Printf("parsing time was not successful - %s", err.Error())
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
