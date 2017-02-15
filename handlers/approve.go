package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func Approve(context echo.Context) error {
	//err := helpers.Approve()
	//if err != nil {
	//return err
	//}

	return context.JSON(http.StatusOK, "Image approved")
}
