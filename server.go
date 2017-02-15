package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/byuoitav/video-capture-and-tweet/handlers"
	"github.com/byuoitav/video-capture-and-tweet/helpers"
	"github.com/byuoitav/video-capture-and-tweet/tweeter"
	"github.com/byuoitav/video-capture-and-tweet/views"
	"github.com/jessemillar/health"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	configptr := flag.String("c", "./config.json", "configuration file")
	prodptr := flag.Bool("p", false, "if is production")

	flag.Parse()

	config := *configptr
	tweeter.Production = *prodptr

	log.Printf("config: %v", config)

	tweeter.Config = tweeter.GetConfiguration(config)
	tweeter.StartChannel = make(chan bool, 1)

	go tweeter.Startup()

	port := ":9000"
	router := echo.New()
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.CORS())

	templateEngine := &helpers.Template{
		Templates: template.Must(template.ParseGlob("public/*/*.html")),
	}

	router.Renderer = templateEngine

	router.GET("/health", echo.WrapHandler(http.HandlerFunc(health.Check)))

	// Views
	router.Static("/*", "public")
	router.GET("/", views.Main)

	router.GET("/tweeter/start", handlers.Start)
	router.GET("/tweeter/stop", handlers.Stop)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
