package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

type Response struct {
	Token string `json:"token"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	client := connectDB()
	db := client.Database(os.Getenv("DB"))
	collection := db.Collection(os.Getenv("DB_COLLECTION"))

	router := http.NewServeMux()

	// Will only accept GET Routes. If not mentioned, route will accept any method
	router.HandleFunc("GET /status", func(w http.ResponseWriter, r *http.Request) {
		var status string = http.StatusText(200)
		fmt.Fprintln(w, status)
	})

	router.HandleFunc("POST /enter", func(w http.ResponseWriter, r *http.Request) {
		reqData, err := readBody(r)
		if err != nil {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		id := insertOne(collection, reqData)
		fmt.Fprintf(w, "Insert ID: %v", id)
	})

	router.HandleFunc("POST /getByEmail", func(w http.ResponseWriter, r *http.Request) {
		reqData, err := readBody(r)
		if err != nil {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

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
	})

	server := http.Server{
		Addr:    ":8080", //Port
		Handler: Logging(router),
	}

	fmt.Println("Server is live on", server.Addr)
	server.ListenAndServe()
}
