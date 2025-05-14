package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type requestStruct struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
	Token string `json:"token"`
	Addr  string `json:"addr"`
}

// TODO:
//	[x] 1. Logging
//	[] 2. Auth

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		message := "From " + r.RemoteAddr + ": " + r.Method + " " + r.URL.Path + " " + time.Now().Format(time.DateTime) + " " + "\n"
		strByte := []byte(message)

		logFile, err := os.OpenFile("api-log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("Could not open log file:", err)
		}
		defer logFile.Close()

		writer := bufio.NewWriter(logFile)
		if _, err := writer.Write(strByte); err != nil {
			panic(err)
		}
		writer.Flush()
	})
}

func readBody(r *http.Request) (requestStruct, error) {
	decoder := json.NewDecoder(r.Body)

	var t requestStruct
	t.Addr = r.RemoteAddr

	if err := decoder.Decode(&t); err != nil {
		return t, err
	}

	fmt.Printf("Token Value: %s \n", t.Token)
	fmt.Printf("Addr Value: %s \n", t.Addr)

	return t, nil

}

func connectDB() *mongo.Client {
	// fmt.Printf("URI: %s", os.Getenv("MONGO_URI"))
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

func insertOne(collection *mongo.Collection, req requestStruct) interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	document := bson.D{
		{Key: "name", Value: req.Name},
		{Key: "email", Value: req.Email},
		{Key: "age", Value: req.Age},
		{Key: "token", Value: req.Token},
		{Key: "addr", Value: req.Addr},
	}
	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted document ID:", result.InsertedID)
	return result.InsertedID
}
func getOneByEmail(collection *mongo.Collection, req requestStruct) (interface{}, error) {

	document := bson.D{{Key: "email", Value: req.Email}}

	var result requestStruct
	err := collection.FindOne(context.TODO(), document).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, err
		}
		panic(err)
	}

	return result, err
}
