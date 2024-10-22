package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func dynamicHandler(w http.ResponseWriter, r *http.Request) {
	// Remove the leading slash and append ".html"
	page := r.URL.Path[1:] // e.g., "/about" -> "about"
	if page == "" {
		page = "index"
	}

	tmplPath := filepath.Join("templates", page+".html")

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
func succesfullPullHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Pulled and restarted successfully"))
	time.Sleep(5 * time.Second)
	http.Redirect(w, r, "/", http.StatusFound)
}

func pull_and_restart(w http.ResponseWriter, r *http.Request) {
	fmt.Println("pulling and rebooting")

	// Initiate the script execution in a goroutine
	scriptPath := "/home/nightflavor/httpserver/pull.sh"
	go func() {
		cmd := exec.Command("bash", "-c", scriptPath)

		// Capture the output of the script
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error running script: %v, output: %s", err, string(output))
			return // Log the error without affecting the HTTP response
		}
		fmt.Println("Script executed successfully:", string(output))
	}()

	// Redirect the user to the success page
	http.Redirect(w, r, "/pullsucces", http.StatusFound)
}

func main() {
	http.HandleFunc("/", dynamicHandler)
	http.HandleFunc("/pull", pull_and_restart)
	http.HandleFunc("/pullsucces", succesfullPullHandler)
	cssFileServer := http.FileServer(http.Dir("./css"))
	http.Handle("/css/", http.StripPrefix("/css/", cssFileServer))
	port := "80"
	fmt.Println("starting server on port", port)
	log.Printf("Starting server on port %s", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
