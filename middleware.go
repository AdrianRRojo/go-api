package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type requestStruct struct {
	Addr, Token string
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

func readBody(r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	for {
		var t requestStruct
		t.Addr = r.RemoteAddr

		if err := decoder.Decode(&t); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Token Value: %s \n", t.Token)
		fmt.Printf("Addr Value: %s \n", t.Addr)
	}
}
