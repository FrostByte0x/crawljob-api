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
		"PURGE_FILES_AGE_IN_HOURS", os.Getenv("PURGE_FILES_AGE_IN_HOURS"),
	)
	// start the purge job, use the same destinationFolder retrieved in the handler package
	enablePurge, err := strconv.ParseBool(os.Getenv("ENABLE_PURGE"))
	if err != nil {
		slog.Warn(err.Error())
	}
	// Default is 24 hours, but can be overridden using the PURGE_FILES_AGE_IN_HOURS.
	if enablePurge {
		purgeMaximumFileAge := 24 // this is the default if there are no values in the container.
		if value, err := strconv.Atoi(os.Getenv("PURGE_FILES_AGE_IN_HOURS")); err == nil {
			purgeMaximumFileAge = value
		}
		jobs.StartPurgeRoutine(handler.GetDestinationFolder(), purgeMaximumFileAge)
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
