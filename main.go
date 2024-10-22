package main

import (
	"fmt"
	"html/template"
	"log"
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
		log.Printf("Template not found: %s, error: %v", tmplPath, err)
		return
	}

	// Define the data to pass to the template
	data := struct {
		Title  string
		Header string
	}{
		Title:  page,
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

func pull_and_restart(w http.ResponseWriter, r *http.Request) {
	fmt.Println("pulling and rebooting")

	// Use the absolute path directly, since we know where the script is located
	scriptPath := "/home/nightflavor/httpserver/pull.sh"
	cmd := exec.Command("bash", "-c", scriptPath)

	// Capture the output of the script
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, "Failed to run script", http.StatusInternalServerError)
		log.Printf("Error running script: %v, output: %s", err, string(output))
		fmt.Println("Error running script:", err)
		fmt.Println("Script output:", string(output))
		return
	}

	// Print the output to both the console and the log
	fmt.Println("Script output:", string(output))
	log.Printf("Script output: %s", string(output))

	// Respond to the HTTP request with the script output
	fmt.Fprintf(w, "Script executed successfully: %s", string(output))
}


func main() {
	http.HandleFunc("/", dynamicHandler)
	http.HandleFunc("/pull", pull_and_restart)
	cssFileServer := http.FileServer(http.Dir("./css"))
	http.Handle("/css/", http.StripPrefix("/css/", cssFileServer))
	port := "80"
	fmt.Println("starting server on port", port)
	log.Printf("Starting server on port %s", port)

	// Use log.Fatal to catch errors with ListenAndServe
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
