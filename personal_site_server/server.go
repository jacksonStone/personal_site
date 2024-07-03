package main

import (
	"embed"
	"log"
	"net/http"
	"path/filepath"
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
		returnStaticHTML(w, "public/jackson-stone.html")
	} else {
		path := r.URL.Path[1:]
		firstPartOfPath := strings.Split(path, "/")[0]
		if firstPartOfPath == "favicon.ico" {
			// todo return favicon file
			data, err := content.ReadFile(filepath.Join("public", "favicon.ico"))
			if err != nil {
				http.Error(w, "Not Found", http.StatusNotFound)
			}
			w.Header().Set("Content-Type", "image/x-icon")
			w.Write(data)
			return
		}
		if firstPartOfPath == ".." || firstPartOfPath == "." {
			http.Error(w, "Invalid Path: Nice Try! No Monkey business!", http.StatusNotFound)
			return
		}
		returnStaticHTML(w, filepath.Join("public", firstPartOfPath+".html"))
	}
}
