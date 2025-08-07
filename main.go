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
	http.HandleFunc("/apply", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.ShowApplyForm(w, r)
		} else if r.Method == http.MethodPost {
			handlers.ApplyHandler(w, r)
		}
})

	log.Println("Server started at http://localhost:9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
	
}