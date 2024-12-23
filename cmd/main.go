package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Start server")

	http.HandleFunc("/", Handler)

	http.ListenAndServe(":8080", nil)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server start port 8080")
}
