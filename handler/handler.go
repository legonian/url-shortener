// Package handler provides functions that represent routing logic
package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/legonian/url-shortener/database"
)

type (
	// Data type that coming from client to create Data
	Url struct {
		Url string `json:"url"`
	}
	// Data represent page content
	Page struct {
		Title      string
		Body       database.Data
		StatusCode int
		StatusText string
	}
)

var templates = template.Must(template.ParseGlob("templates/*"))

// Send Index page
func Index(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", Page{Title: "Home Page"})
}

// Send Info page about given shortcut from database (or cache) to copy and
// view clicks count of that shortcut
func Info(w http.ResponseWriter, r *http.Request) {
	shortcut := chi.URLParam(r, "shortcut")
	cacheData := database.CheckCache(shortcut, database.NotViewed)
	if cacheData.OK {
		renderTemplate(w, "info", Page{
			Title: "URL Info",
			Body:  cacheData,
		})
		return
	}
	data := database.GetData(shortcut, database.NotViewed)
	if !data.OK {
		renderTemplate(w, "error", Page{
			Title:      "Not found in database",
			StatusCode: http.StatusNotFound,
		})
		return
	}
	database.AddCache(data)
	renderTemplate(w, "info", Page{
		Title: "URL Info",
		Body:  data,
	})
}

// Parse POST request, send given URL to database, after getting data with
// shortcut code send it to client in json
func SetRedirectJson(w http.ResponseWriter, r *http.Request) {
	var postData Url
	err := json.NewDecoder(r.Body).Decode(&postData)
	if err != nil {
		sendJson(w, database.Data{OK: false}, http.StatusBadRequest)
		return
	}
	newUrl := string(postData.Url)
	if !isValidUrl(newUrl) {
		sendJson(w, database.Data{OK: false}, http.StatusBadRequest)
		return
	}
	res := database.CreateData(newUrl)
	sendJson(w, res, http.StatusCreated)
}

// Parse GET request, and after getting data about given given shortcut
// link from database (or cache) redirect client to full link
func Redirect(w http.ResponseWriter, r *http.Request) {
	shortcut := chi.URLParam(r, "shortcut")
	cacheData := database.CheckCache(shortcut, database.IsViewed)
	if cacheData.OK {
		http.Redirect(w, r, cacheData.FullURL, http.StatusFound)
		return
	}
	data := database.GetData(shortcut, database.IsViewed)
	if !data.OK {
		renderTemplate(w, "error", Page{
			Title:      "Not found in database",
			StatusCode: http.StatusNotFound,
		})
		return
	}
	database.AddCache(data)
	http.Redirect(w, r, data.FullURL, http.StatusFound)
}

// Check does given string is valid URL
func isValidUrl(urlToCheck string) bool {
	_, err := url.ParseRequestURI(urlToCheck)
	if err != nil {
		log.Printf("Bad URL: %s", urlToCheck)
		return false
	}
	return true
}

// Send given template file or error page
func renderTemplate(w http.ResponseWriter, tmpl string, p Page) {
	if p.StatusCode == 0 {
		p.StatusCode = http.StatusOK
	}
	w.WriteHeader(p.StatusCode)
	if p.StatusText == "" {
		p.StatusText = http.StatusText(p.StatusCode)
	}
	if p.Title == "" {
		p.Title = http.StatusText(p.StatusCode)
	}
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Send given data struck to writer as json respond
func sendJson(w http.ResponseWriter, dataToSend database.Data, httpCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	err := json.NewEncoder(w).Encode(dataToSend)
	if err != nil {
		renderTemplate(w, "error", Page{
			Title:      "json error",
			StatusCode: http.StatusNotFound,
		})
	}
}
