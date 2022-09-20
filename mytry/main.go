package main

import (
	"fmt"
	"log"
	"net/http"
)

func messageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is my try to up the server")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/mytry", messageHandler)
	log.Println("Listening...")

	http.ListenAndServe(":8080", mux)
}
