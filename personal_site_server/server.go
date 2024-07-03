package main

import (
	"embed"
	"log"
	"net/http"
	"strings"
)

//go:embed public
var content embed.FS

func main() {

	http.HandleFunc("/", rootHandler)
	// For local development, just use http
	log.Println("Starting server on :3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}

func returnStaticHTML(w http.ResponseWriter, path string) {
	data, err := content.ReadFile(path)
	if err != nil {
		http.Error(w, "Invalid Path: This page does not exist.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		returnStaticHTML(w, "public/about.html")
	} else {
		path := r.URL.Path[1:]
		firstPartOfPath := strings.Split(path, "/")[0]
		returnStaticHTML(w, "public/"+firstPartOfPath+".html")
	}
}
