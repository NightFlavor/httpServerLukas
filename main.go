package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		http.Error(w, "Failed to find home directory", http.StatusInternalServerError)
		return
	}
	scriptPath := filepath.Join(homeDir, "httpserver/pull.sh")

	cmd := exec.Command("bash", "-c", scriptPath)
	output, err := cmd.CombinedOutput() // CombinedOutput captures both stdout and stderr
	if err != nil {
		// Print detailed error information
		fmt.Println("Error running script:", err)
		fmt.Println("Script output:", string(output))
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
