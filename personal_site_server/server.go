package main

import (
	"embed"
	"log"
	"net/http"
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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	data, err := content.ReadFile("public/about.html")
	if err != nil {
		http.Error(w, "Could not read embedded file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}
