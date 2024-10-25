package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func dynamicHandler(w http.ResponseWriter, r *http.Request) {
	// Remove the leading slash and append ".html"
	page := r.URL.Path[1:] // e.g., "/about" -> "about"
	pageTitle := "Title"
	pageName := page
	if page == "" {
		page = "index"
		pageName = "Home"
		pageTitle = "Welcome to Lukas' portfolio:"
	}

	cookie, err := r.Cookie("mode")
	mode := "light" // Default mode
	if err == nil {
		mode = cookie.Value
	}

	cssFile := "light.css"
	if mode == "dark" {
		cssFile = "dark.css"
	}

	tmplPath := filepath.Join("templates", page+".html")

	if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
		http.Error(w, "404 Not Found in dynamicHandler", http.StatusNotFound)
		log.Printf("Template not found: %s, error: %v", tmplPath, err)
		return
	}
	fmt.Print("/" + r.URL.Path)
	// Define the data to pass to the template
	data := struct {
		Title     string
		CSS       string
		Redirect  string
		PageTitle string
	}{
		Title:     pageName,
		CSS:       cssFile,
		Redirect:  r.URL.Path,
		PageTitle: pageTitle,
	}

	// Render the template if it exists
	renderTemplate(w, tmplPath, data)
}

func setModeHandler(w http.ResponseWriter, r *http.Request) {
	// Get the mode from the query string (e.g., ?mode=dark)
	mode := r.URL.Query().Get("mode")
	redirect := r.URL.Query().Get("redirect")
	if mode != "light" && mode != "dark" {
		mode = "light" // Default to light if invalid value
	}

	// Set the mode cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "mode",
		Value:   mode,
		Expires: time.Now().Add(365 * 24 * time.Hour), // 1-year expiration
		Path:    "/",
	})
	// Redirect back to the homepage (or wherever you want)
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// Parse layout and requested template
	t, err := template.ParseFiles("templates/layout.html", tmpl)
	if err != nil {
		http.Error(w, "404 Not Found in renderTemplate", http.StatusNotFound)
		log.Printf("Failed to parse templates: %v, error: %v", tmpl, err)
		return
	}

	// Execute the template with the provided data
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error executing template %s: %v", tmpl, err)
	}
}

func main() {
	http.HandleFunc("/", dynamicHandler)

	cssFileServer := http.FileServer(http.Dir("./css"))

	http.HandleFunc("/set-mode", setModeHandler)
	http.Handle("/css/", http.StripPrefix("/css/", cssFileServer))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := "8000"
	fmt.Println("starting server on port", port)
	log.Printf("Starting server on port %s", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
