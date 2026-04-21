package handler

import (
	"crawljob-api/model"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

type ValidBody struct {
	Url string `json:"url"`
}

// The handler is the code that handles each http request received by the server

func Handle(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed) // 405
		slog.Warn("Method not allowed", "method", request.Method, "remote", request.RemoteAddr)
		fmt.Fprintf(writer, "method not allowed")
		return
	}
	// this is an empty var of validBody type.
	var requestUrl ValidBody
	// read the stream
	requestBody, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest) // 400
		slog.Info("Received wrong body from a client.")
		fmt.Fprintf(writer, "Invalid body: %s", err)
		return
	}
	// unmarshal the json into the var created above
	err = json.Unmarshal(requestBody, &requestUrl)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest) // 400
		slog.Info("Received wrong body from a client.")
		fmt.Fprintf(writer, "Invalid body: %s", err)
		return
	}
	// Now requestUrl.Url is our request URL!
	err = ValidateUrl(requestUrl.Url)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		slog.Info("Received invalid URL from a client", "url", requestUrl.Url)
		fmt.Fprintf(writer, "Invalid URL received %s", err)
		return
	}
	// The following code is the client return + http/201
	slog.Info("Received a valid request from", "IP Addr", request.RemoteAddr, "url", requestUrl.Url)

	// Start of main processing for jobs
	// Check destination for the job file
	destinationFileJob := os.Getenv("CRAWLJOB_FOLDER")
	if destinationFileJob == "" {
		slog.Warn("CRAWLJOB_FOLDER is not set, using fallback to current directory")
		destinationFileJob = "." // current running directory
	} else if _, err := os.Stat(destinationFileJob); os.IsNotExist(err) {
		slog.Warn("CRAWLJOB_FOLDER does not exist, using fallback to current directory")
		destinationFileJob = "." // current running directory
	}
	// Check the destination for the download directory
	destinationFolder := os.Getenv("DESTINATION_FOLDER")
	if destinationFolder == "" {
		slog.Warn("DESTINATION_FOLDER is not set")
	}
	// Generate crawl job file
	err = model.GenerateJobFile(requestUrl.Url, destinationFolder, destinationFileJob)
	if err != nil {
		slog.Info("Error generating crawljob file", "error", err)
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "Error handling job file: %s", err)
		return
	}
	// finish with this
	writer.WriteHeader(http.StatusCreated) // 201 is the expected request for resource creation in Rest APIs.
	fmt.Fprintf(writer, "Job successfully received")
}
