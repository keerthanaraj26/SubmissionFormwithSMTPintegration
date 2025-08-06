package main

import (
	"log"
	"net/http"
	"studentform/handlers"
)

func main() {
	http.HandleFunc("/", handlers.FormHandler)
	http.HandleFunc("/submit", handlers.SubmitHandler)

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}