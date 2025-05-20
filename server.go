package main

import (
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func runServer(collection *mongo.Collection) {
	router := http.NewServeMux()

	// Will only accept GET Routes. If not mentioned, route will accept any method
	router.HandleFunc("GET /status", func(w http.ResponseWriter, r *http.Request) {
		var status string = http.StatusText(200)
		fmt.Fprintln(w, status)
	})

	router.HandleFunc("POST /enter", func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("Handler invoked: Method=%s, Path=%s, RemoteAddr=%s\n", r.Method, r.URL.Path, r.RemoteAddr)
		reqData, err := readBody(r)
		if err != nil {
			http.Error(w, "Invalid Request ", http.StatusBadRequest)
			return
		}
		// fmt.Printf("Token Value (raw): %q\n", reqData.Token)

		companyID, isAuth := Auth(reqData)
		if isAuth {
			_, error := insertOne(collection, reqData, companyID)

			if error != nil {
				http.Error(w, "Error: Could not submit request", http.StatusBadRequest)
			}
			// fmt.Fprintf(w, "Insert ID: %v", id)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "Successfully submitted request")

		} else {
			http.Error(w, "Invalid Token", http.StatusBadRequest)
			return
		}
	})

	router.HandleFunc("POST /getByEmail", func(w http.ResponseWriter, r *http.Request) {
		reqData, err := readBody(r)
		if err != nil {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		_, isAuth := Auth(reqData)
		if isAuth {

			user, err := getOneByEmail(collection, reqData)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					http.Error(w, "User not found", http.StatusNotFound)
					return
				}
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "User: %+v", user)
		} else {
			http.Error(w, "Invalid Token", http.StatusBadRequest)
			return
		}

	})

	server := http.Server{
		Addr:    ":8080", //Port
		Handler: Logging(router),
	}

	fmt.Println("Server is live on", server.Addr)
	server.ListenAndServe()
}
