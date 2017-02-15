package views

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
)

func Main(context echo.Context) error {
	log.Println("Returning main page")
	return context.Render(http.StatusOK, "main", "")
}

func Toggle(context echo.Context) error {
	log.Println("Returning toggle page")
	return context.Render(http.StatusOK, "toggle", "")
}

func Approve(context echo.Context) error {
	log.Println("Returning approval page")
	return context.Render(http.StatusOK, "approve", "")
}
