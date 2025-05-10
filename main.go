package main

import (
	"fmt"
	"net/http"
)

type Response struct {
	Token string `json:"token"`
}

func main() {
	router := http.NewServeMux()

	// Will only accept GET Routes. If not mentioned, route will accept any method
	router.HandleFunc("GET /status", func(w http.ResponseWriter, r *http.Request) {
		var status string = http.StatusText(200)
		fmt.Fprintln(w, status)
	})

	router.HandleFunc("POST /enter", func(w http.ResponseWriter, r *http.Request) {
		readBody(r)
	})

	router.HandleFunc("GET /getStatus/{id...}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintln(w, id)
	})

	server := http.Server{
		Addr:    ":8080", //Port
		Handler: Logging(router),
	}

	fmt.Println("Server is live on", server.Addr)
	server.ListenAndServe()
}
