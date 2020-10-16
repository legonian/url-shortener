package main

import (
	"log"
	"net/http"
	"os"
	
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.Default()

	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", indexRouter)
	router.POST("/submit", submitRouter)

	router.Run(":" + port)
}

func indexRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl.html", nil)
}

func submitRouter(c *gin.Context) {
	log.Printf(c.PostForm("url"))
	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"generatedURL": "qwe",
	})
}
