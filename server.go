package main

import (
	"html/template"
	"net/http"

	"github.com/jessemillar/health"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	port := ":9000"
	router := echo.New()
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.CORS())

	templateEngine := &helpers.Template{
		Templates: template.Must(template.ParseGlob("public/*/*.html")),
	}

	router.Renderer = templateEngine

	router.GET("/health", echo.WrapHandler(http.HandlerFunc(health.Check)))
	router.GET("/approve", handlerGroup.Approve)

	// Views
	router.Static("/*", "public")
	router.GET("/", views.Main)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
