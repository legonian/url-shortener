// Package main provides Echo framework initialization and set handlers to
// their path
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/legonian/url-shortener/database"
	"github.com/legonian/url-shortener/handler"
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
	app := SetupApp()
	actionOnInterrupt()
	log.Fatal(http.ListenAndServe(":3000", app))
}

func SetupApp() chi.Router {
	// Initialize database
	err := database.Init()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", handler.Index)
	r.Get("/{shortcut}", handler.Redirect)
	r.Get("/{shortcut}/info", handler.Info)

	r.Post("/create", handler.SetRedirectJson)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "public"))
	FileServer(r, "/public", filesDir)
	return r
}

// func main() {
// 	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
// 		Format: "path:${uri} | ${method} method to ${status} | t=${latency_human}\n",
// 	}))
// 	e.Use(middleware.Recover())
// 	e.Use(middleware.Secure())
// }

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

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
