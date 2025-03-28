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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var db *sql.DB

func connectDb() {
	connStr := "host=localhost port=5432 user=lukas password=admin dbname=lukasweb sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr) // Assign to global db
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Connection dead:", err)
	}
	fmt.Println("Connected!")
}

var (
	titleCaser = cases.Title(language.English)
)

type Card struct {
	ID          int
	Section     string
	Title       string
	ImageLink   sql.NullString
	LinkURL     sql.NullString
	Description sql.NullString
}

func getCardsBySection() (map[string][]Card, error) {
	rows, err := db.Query(`
        SELECT section, title, image_link, link_url, description 
        FROM cards 
        ORDER BY section, id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sections := make(map[string][]Card)
	for rows.Next() {
		var c Card
		err := rows.Scan(
			&c.Section,
			&c.Title,
			&c.ImageLink,
			&c.LinkURL,
			&c.Description,
		)
		if err != nil {
			return nil, err
		}
		sections[c.Section] = append(sections[c.Section], c)
	}
	return sections, nil
}

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

	log.Println("Looking for template at:", tmplPath)
	if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
		log.Printf("Template not found: %s", tmplPath)
	}
	if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
		http.Error(w, "404 Not Found in dynamicHandler", http.StatusNotFound)
		log.Printf("Template not found: %s, error: %v", tmplPath, err)
		return
	}
	sections, err := getCardsBySection()
	if err != nil {
		log.Printf("Failed to get cards: %v", err)
	}

	data := struct {
		Title     string
		CSS       string
		Redirect  string
		PageTitle string
		Sections  map[string][]Card
	}{
		Title:     pageName,
		CSS:       cssFile,
		Redirect:  r.URL.Path,
		PageTitle: pageTitle,
		Sections:  sections,
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

func renderTemplate(w http.ResponseWriter, tmpl string, data any) {
	// Create a new template set that includes both layout and content
	tmplSet := template.Must(template.New("").Funcs(template.FuncMap{
		"title":     titleCaser.String,
		"lower":     strings.ToLower,
		"hasPrefix": strings.HasPrefix,
		"where": func(items map[string][]Card, section string) []Card {
			return items[section]
		},
	}).ParseFiles("templates/layout.html", tmpl))

	// Execute the layout template which will include the content
	err := tmplSet.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}
func main() {
	// Set up handlers
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/", dynamicHandler)
	http.HandleFunc("/set-mode", setModeHandler)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir("./media"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Database connection
	connectDb()
	defer db.Close()

	// Start server
	port := "8000"
	fmt.Println("Starting server on port:", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
