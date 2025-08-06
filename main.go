package main

import (
	"log"
	"net/http"
	"studentform/handlers"
)

func main() {
	http.HandleFunc("/", handlers.RegisterHandler)
	http.HandleFunc("/submit", handlers.SubmitHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	log.Println("Server started at http://localhost:9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
	
}