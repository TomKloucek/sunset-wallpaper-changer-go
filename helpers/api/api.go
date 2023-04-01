package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	timeHelper "sunset-wallpaper-changer-go/helpers/time"
	"sunset-wallpaper-changer-go/logger"
	"sunset-wallpaper-changer-go/objects"
)

var LOGGER = logger.GetInstance().Logger

// GetResponseFromSunsetApi gets response from api based on user defined latitude and longitude in type Response
func GetResponseFromSunsetApi(latStr string, longStr string, offsetHours int) objects.Response {
	apiUrl := fmt.Sprintf("https://api.sunrise-sunset.org/json?lat=%s&lng=%s", latStr, longStr)

	response, err := http.Get(apiUrl)
	if err != nil {
		LOGGER.Printf("%s - fetching api was not successful", err.Error())
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		LOGGER.Printf("%s - reading response was not successful", err.Error())
	}

	var parsedResponse objects.Response
	err = json.Unmarshal([]byte(string(responseData)), &parsedResponse)
	if err != nil {
		LOGGER.Printf("%s - parsing response was not successful", err.Error())
		panic(err)
	}

	parsedResponse.Results.Sunrise, err = timeHelper.AddHoursToTimeString(parsedResponse.Results.Sunrise, offsetHours)
	if err != nil {
		LOGGER.Printf("%s - adding time to sunrise was not successful", err.Error())
		panic(err)
	}
	parsedResponse.Results.Sunset, err = timeHelper.AddHoursToTimeString(parsedResponse.Results.Sunset, offsetHours)
	if err != nil {
		LOGGER.Printf("%s - adding time to sunset was not successful", err.Error())
		panic(err)
	}
	return parsedResponse
}
