package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func give_website(w http.ResponseWriter, r *http.Request) {
	htmlFile, _ := os.Open("html/index.html")
	defer htmlFile.Close()
	w.Header().Set("Content-Type", "text/html")
	io.Copy(w, htmlFile)
}
func pull_and_restart(w http.ResponseWriter, r *http.Request) {
	fmt.Println("pulling and rebooting")
	cmd := exec.Command("pull.sh")
	output, _ := cmd.Output()
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
