package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}
func dynamicHandler(w http.ResponseWriter, r *http.Request) {
	// Remove the leading slash and append ".html"
	page := r.URL.Path[1:] // e.g., "/about" -> "about"
	pageTitle := "Title"
	pageName := page
	if page == "" {
		page = "index"
		pageName = "Home"
		pageTitle = "Lukas' portfolio:"
	} else {
		pageTitle = "Lukas' " + page + " page"
	}

	w.Header().Set("Accept-CH", "Sec-CH-Prefers-Color-Scheme")
	mode := r.Header.Get("Sec-CH-Prefers-Color-Scheme")
	mode = If(mode != "", mode, "dark") //ternary terrorism

	cookie, err := r.Cookie("mode")
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

	renderTemplate(w, tmplPath, data)
}

func setModeHandler(w http.ResponseWriter, r *http.Request) {
	//(?mode=dark)
	mode := r.URL.Query().Get("mode")
	redirect := r.URL.Query().Get("redirect")
	if mode != "light" && mode != "dark" {
		mode = "light"
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "mode",
		Value:   mode,
		Expires: time.Now().Add(365 * 24 * time.Hour), // 1y
	})
	//redirect back
	if strings.Contains(redirect, ".") {
		return
	}
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("templates/layout.html", tmpl)
	if err != nil {
		http.Error(w, "404 Not Found in renderTemplate", http.StatusNotFound)
		log.Printf("Failed to parse templates: %v, error: %v", tmpl, err)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error executing template %s: %v", tmpl, err)
	}
}

func connectDb() {
	connStr := "postgres://postgres:admin@localhost:5432/gopgtest?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Error pinging db: %v", err)
	}
	defer db.Close()
}

func main() {
	http.HandleFunc("/", dynamicHandler)
	http.HandleFunc("/set-mode", setModeHandler)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir("./media"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	connectDb()

	port := "8000"
	fmt.Println("Starting server on port: ", port)
	fmt.Println("Succes!")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
