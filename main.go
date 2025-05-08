package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		var status string = http.StatusText(200)
		fmt.Fprintln(w, status)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("SERVER ERROR: ", err)
	}
}
