package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
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
	mode := flag.String("mode", "server", "Run Mode: 'server' or 'jwt")

	flag.Parse()

	switch *mode {
	case "jwt":
		generateJWT(collection)
	default:
		runServer(collection)
	}

}
