package main

import (
	"io"
	"net/http"
	"os"
)

func give_website(w http.ResponseWriter, r *http.Request) {
	htmlFile, _ := os.Open("html/index.html")
	defer htmlFile.Close()
	w.Header().Set("Content-Type", "text/html")
	io.Copy(w, htmlFile)
}

func main() {
	http.HandleFunc("/index.html", give_website)

	fileServer := http.FileServer(http.Dir("./html"))
	http.Handle("/html/", http.StripPrefix("/html/", fileServer))

	http.ListenAndServe(":80", nil)
}
