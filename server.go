// Package main provides Echo framework initialization and set handlers to
// their path
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	catchExit()
	log.Fatal(http.ListenAndServe(":"+port, app))
}

// Initialize database and router
func SetupApp() chi.Router {
	if err := database.Init(); err != nil {
		log.Fatal(err)
	}

	if err := handler.SetTemplates("templates/*"); err != nil {
		log.Fatal(err)
	}

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

// Handle exit signals from OS and do action to to prevent losing data
func catchExit() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		database.ClearCache()
		fmt.Println("Cache saved before exit")
		os.Exit(0)
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
