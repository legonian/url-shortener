// Package main provides Echo framework initialization and set handlers to
// their path
package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/legonian/url-shortener/database"
	"github.com/legonian/url-shortener/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var port string

func init() {
	if os.Getenv("GO_ENABLE_LOG") == "" {
		log.SetOutput(ioutil.Discard)
	}
	// Check PORT env variable
	port = os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
}

func main() {
	// Initialize database
	err := database.Init()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize app
	e := echo.New()

	// Middlewares
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "path:${uri} | ${method} method to ${status} | t=${latency_human}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	// Static Files including HTML
	e.Static("/public", "public")

	// Main Routes
	e.GET("/", handler.Index)
	e.GET("/:short_url", handler.Redirect)
	e.GET("/:short_url/info", handler.Info)

	// Called from Index Page javascript
	e.POST("/create", handler.SetRedirectJson)
	// Called from Info Page javascript
	e.POST("/:short_url/json", handler.InfoJson)

	actionOnInterrupt()

	e.Logger.Fatal(e.Start(":" + port))
}

func actionOnInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			database.ClearCache()
			log.Fatal(sig)
		}
	}()
}
