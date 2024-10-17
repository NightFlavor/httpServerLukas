package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func give_website(w http.ResponseWriter, r *http.Request) {
	htmlFile, err := os.Open("html/index.html")
	if err != nil {
		http.Error(w, "Could not open file", http.StatusInternalServerError)
		return
	}
	defer htmlFile.Close()
	w.Header().Set("Content-Type", "text/html")
	io.Copy(w, htmlFile)
}

func pull_and_restart(w http.ResponseWriter, r *http.Request) {
	fmt.Println("pulling and rebooting")
	cmd := exec.Command("~httpserver/pull.sh")
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, "Failed to run script", http.StatusInternalServerError)
		return
	}
	fmt.Println(string(output))
}

func main() {
	http.HandleFunc("/", give_website)
	fileServer := http.FileServer(http.Dir("./html"))
	http.HandleFunc("/pull", pull_and_restart)
	http.Handle("/html/", http.StripPrefix("/html/", fileServer))
	port := "12350"
	fmt.Println("starting server on ", port)
	http.ListenAndServe(":"+port, nil)
}
