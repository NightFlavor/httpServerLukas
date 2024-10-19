package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func dynamicHandler(w http.ResponseWriter, r *http.Request) {
	// Remove the leading slash and append ".html"
	page := r.URL.Path[1:] // e.g., "/about" -> "about"
	if page == "" {
		page = "index" // Default to index if no page is specified
	}

	// Construct the template path
	tmplPath := filepath.Join("templates", page+".html")

	// Check if the requested template exists
	if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
		http.Error(w, "404 Not Found in dynamicHandler", http.StatusNotFound)
		return
	}

	// Define the data to pass to the template
	data := struct {
		Title  string
		Header string
	}{
		Title:  page, // You can customize this or make it dynamic
		Header: "Welcome to " + page,
	}

	// Render the template if it exists
	renderTemplate(w, tmplPath, data)
}
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// Parse layout and requested template
	t, err := template.ParseFiles("templates/layout.html", tmpl)
	if err != nil {
		http.Error(w, "404 Not Found in renderTemplate", http.StatusNotFound)
		return
	}

	// Execute the template with the provided data
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
func pull_and_restart(w http.ResponseWriter, r *http.Request) {
	fmt.Println("pulling and rebooting")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		http.Error(w, "Failed to find home directory", http.StatusInternalServerError)
		return
	}
	scriptPath := filepath.Join(homeDir, "httpserver/pull.sh")
	cmd := exec.Command("bash", "-c", scriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running script:", err)
		fmt.Println("Script output:", string(output))
		http.Error(w, "Failed to run script", http.StatusInternalServerError)
		return
	}
	fmt.Println(string(output))
}

func main() {
	http.HandleFunc("/", dynamicHandler)
	http.HandleFunc("/pull", pull_and_restart)
	cssFileServer := http.FileServer(http.Dir("./css"))
	http.Handle("/css/", http.StripPrefix("/css/", cssFileServer))
	port := "80"
	fmt.Println("starting server on ", port)
	http.ListenAndServe(":"+port, nil)
}
