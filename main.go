package main

import (
	"fmt"
	"io"
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
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Could not read request body", http.StatusBadRequest)
			fmt.Println("Error reading request body", err)
		}

		defer r.Body.Close()

		fmt.Fprintf(w, "Request body: %s \n", string(body))
	})

	server := http.Server{
		Addr:    ":8080", //Port
		Handler: Logging(router),
	}

	fmt.Println("Server is live on", server.Addr)
	server.ListenAndServe()
}
