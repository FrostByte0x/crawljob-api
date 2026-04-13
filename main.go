package main

import (
	"crawljob-api/handler"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	slog.Info("Starting crawljob-api",
		"DESTINATION_FOLDER", os.Getenv("DESTINATION_FOLDER"),
		"CRAWLJOB_FOLDER", os.Getenv("CRAWLJOB_FOLDER"),
		"ALLOWED_DOMAINS", os.Getenv("ALLOWED_DOMAINS"),
	)

	http.HandleFunc("/", handler.HandleUI)
	http.HandleFunc("/jobs", handler.Handle)
	http.HandleFunc("/api/files", handler.HandleFiles)
	http.HandleFunc("/downloads", handler.HandleDownloadUi)
	http.HandleFunc("/download", handler.DownloadFiles)
	http.HandleFunc("/download/folder", handler.DownloadFolder)

	slog.Info("Server listening", "addr", ":8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		slog.Error("Server failed", "err", err)
	}
}
