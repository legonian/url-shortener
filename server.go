package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/legonian/url-shortener/database"
	"github.com/legonian/url-shortener/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	log.SetOutput(ioutil.Discard)

	// Initialize database
	db, err := database.Init()
	if err != nil {
		log.Fatal(err)
	}

	// Check PORT env variable
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// Initialize app
	e := echo.New()

	// Initialize handler
	h := &handler.Handler{DB: db}

	// Middlewares
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "path:${uri} | ${method} method to ${status} | t=${latency_human}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	// Static Files including HTML
	e.Static("/public", "public")

	// Main Routes
	e.GET("/", h.Index)
	e.GET("/:short_url", h.Redirect)
	e.GET("/:short_url/info", h.Info)

	// Called from Index Page javascript
	e.POST("/create", h.SetRedirectJson)
	// Called from Info Page javascript
	e.POST("/:short_url/json", h.InfoJson)

	e.Logger.Fatal(e.Start(":" + port))
}
