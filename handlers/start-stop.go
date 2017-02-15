package handlers

import (
	"net/http"
	"time"

	"github.com/byuoitav/video-capture-and-tweet/tweeter"
	"github.com/labstack/echo"
)

//Start starts the tweeting
func Start(context echo.Context) error {
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
