package main

import (
	"crawljob-api/handler"
	"crawljob-api/jobs"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

func main() {
	slog.Info("Starting crawljob-api",
		"DESTINATION_FOLDER", os.Getenv("DESTINATION_FOLDER"),
		"CRAWLJOB_FOLDER", os.Getenv("CRAWLJOB_FOLDER"),
		"ALLOWED_DOMAINS", os.Getenv("ALLOWED_DOMAINS"),
		"ENABLE_PURGE", os.Getenv("ENABLE_PURGE"),
	)
	// start the purge job, use the same destinationFolder retrieved in the handler package
	enablePurge, err := strconv.ParseBool(os.Getenv("ENABLE_PURGE"))
	if err != nil {
		slog.Warn(err.Error())
	}
	if enablePurge {
		jobs.StartPurgeRoutine(handler.GetDestinationFolder())
	}
	// register handlers
	http.HandleFunc("/", handler.HandleUI)
	http.HandleFunc("/jobs", handler.Handle)
	http.HandleFunc("/api/files", handler.HandleFiles)
	http.HandleFunc("/downloads", handler.HandleDownloadUi)
	http.HandleFunc("/download", handler.DownloadFiles)
	http.HandleFunc("/download/folder", handler.DownloadFolder)
	// start the web server
	slog.Info("Server listening", "addr", ":8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		slog.Error("Server failed", "err", err)
	}
}
