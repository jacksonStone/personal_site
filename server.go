package main

import (
	"fmt"
	"log"
	"net/http"
)

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
	fmt.Fprintf(w, "Welcome to my personal site!")
}
