package main

import (
	"crawljob-api/handler"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/jobs", handler.Handle)
	fmt.Println("Starting web server on :8080")
	http.ListenAndServe(":8080", nil)
}
