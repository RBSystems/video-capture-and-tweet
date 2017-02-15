package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/byuoitav/video-capture-and-tweet/tweeter"
	"github.com/labstack/echo"
)

//Start starts the tweeting
func Start(context echo.Context) error {
	var err error
	if len(context.QueryParam("interval")) != 0 {
		tweeter.Interval, err = strconv.Atoi(context.QueryParam("interval"))
		if err != nil {
			return context.JSON(http.StatusBadRequest, "invalid interval")
		}
	} else {
		tweeter.Interval = 300 //default interval is 5 minutes
	}
	if tweeter.Interval < 5 {
		return context.JSON(http.StatusBadRequest, "Invalid interval, interval must be 5  or more seconds.")
	}

	tweeter.StartChannel <- true
	return context.JSON(http.StatusOK, "Tweeting started")
}

//Stop stops the tweeting from happening
func Stop(context echo.Context) error {
	tweeter.EndChannel <- true
	select {
	case <-tweeter.ConfirmStop:
		return context.JSON(http.StatusOK, "Process stopped")
	case <-time.After(time.Second * 5):
		return context.JSON(http.StatusInternalServerError, "Could not receive stop signal from tweeter.")
	}
}
