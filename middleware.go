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

	"github.com/golang-jwt/jwt/v5"
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
type tokenStruct struct {
	CompanyID string `json:"companyID"`
	Email     string `json:"email"`
	Exp       int64  `json:"exp"`
	Token     string `json:"token"`
}
type responseStruct struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (r *responseStruct) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
func (r *responseStruct) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return r.ResponseWriter.Write(b)
}

// TODO:
//	[x] 1. Logging
//	[x] 2. Auth

func Logging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &responseStruct{ResponseWriter: w, statusCode: http.StatusOK}

		handler.ServeHTTP(recorder, r)

		req, err := readBody(r)

		var errMessage string
		if err != nil {
			// http.Error(w, "Invalid Request", http.StatusBadRequest)
			fmt.Println("Invalid Request")
			errMessage = err.Error()
		} else {
			errMessage = "No errors reading request"
		}

		companyID, isAuth := Auth(req)

		bodyBytes, _ := json.Marshal(req)

		// message = "From " + r.RemoteAddr + ": \n\t Method:" + r.Method + " Path: " + r.URL.Path + " Time: " + time.Now().Format(time.DateTime) + " \n\t Token Data:" + companyID + " " + "Is Auth: " + fmt.Sprintf("%v", isAuth) + " \n\t Body: " + string(bodyBytes) + " \n\t Error: " + message
		message := fmt.Sprintf(
			"\n From %s:\n\tMethod: %s Path: %s Time: %s\n\tToken Data: CompanyID: %s Is Auth: %v\n\tRequest Body: %s\n\tResponse Status: %d\n\tResponse Body: %s\n\tErrors Reading Request: %s\n\t",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			time.Now().Format(time.DateTime),
			companyID,
			isAuth,
			string(bodyBytes),
			recorder.statusCode,
			string(recorder.body),
			errMessage,
		)

		strByte := []byte(message)

		logFile, err := os.OpenFile("api-log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Could not open log file:", err)
		}
		defer logFile.Close()

		writer := bufio.NewWriter(logFile)
		if _, err := writer.Write(strByte); err != nil {
			log.Println(err)
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
		log.Println(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

func Auth(req requestStruct) (companyID string, isAuth bool) {
	token, err := jwt.Parse(req.Token, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err != nil {
		log.Println("Error reading Token: ", err)
		return "", false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if companyID, ok := claims["companyID"].(string); ok && companyID != "" {
			// check if the token has expired
			if exp, ok := claims["exp"].(float64); ok {
				expirationTime := time.Unix(int64(exp), 0)
				if time.Now().Before(expirationTime) {
					return companyID, true
				}
			}
		}
	}
	return "", false
}

func insertOne(collection *mongo.Collection, req requestStruct, companyID string) interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	document := bson.D{
		{Key: "name", Value: req.Name},
		{Key: "email", Value: req.Email},
		{Key: "age", Value: req.Age},
		{Key: "submittedBy", Value: companyID},
		{Key: "addr", Value: req.Addr},
	}
	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Inserted document ID:", result.InsertedID)
	return result.InsertedID
}

func insertNewCompany(collection *mongo.Collection, t tokenStruct) interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	document := bson.D{
		{Key: "companyID", Value: t.CompanyID},
		{Key: "email", Value: t.Email},
		{Key: "exp", Value: t.Exp},
		{Key: "token", Value: t.Token},
	}
	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
	}

	return result, err
}
