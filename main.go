package main

import (
	"crawljob-api/handler"
	"fmt"
	"net/http"
)

func main() {
	// Handle the main entry point
	http.HandleFunc("/", handler.HandleUI)
	// Handle /jobs
	http.HandleFunc("/jobs", handler.Handle)
	// handle the retrieval of files
	http.HandleFunc("/api/files", handler.HandleFiles)
	// serve the file interface
	http.HandleFunc("/downloads", handler.HandleDownloadUi)
	// Download files
	http.HandleFunc("/download", handler.DownloadFiles)
	// download folders
	http.HandleFunc("/download/folder", handler.DownloadFolder)
	fmt.Println("Starting web server on :8080")
	http.ListenAndServe(":8080", nil)
}
