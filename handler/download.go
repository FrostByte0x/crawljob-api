package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type file struct {
	Name      string // file name, no extension
	Type      string // directory or file
	Extension string // used to display the extension: rar, mkv, etc.
	Size      string // file size
}

func ReturnExtension(s os.DirEntry) string {
	if s.IsDir() {
		return "DIR"
	} else if filepath.Ext(s.Name()) != "" {
		return strings.ToUpper(filepath.Ext(s.Name()))
	} else {
		return "No file extension found"
	}
}

// FormatFileSize returns a neatly formatted string for file size
// 1.2 MB
// 2.4 GB
// 0.2 KB and so on.
func FormatFileSize(size int64) string {
	switch {
	case size < 1024:
		return fmt.Sprintf("%d B", size)
	case size < 1024*1024:
		return fmt.Sprintf("%.1f KB", float64(size)/1024)
	case size < 1024*1024*1024:
		return fmt.Sprintf("%.1f MB", float64(size)/1024/1024)
	case size < 1024*1024*1024*1024:
		return fmt.Sprintf("%.1f GB", float64(size)/1024/1024/1024)
	default:
		return fmt.Sprintf("%.1f TB", float64(size)/1024/1024/1024/1024)
	}
}
func HandleFiles(rw http.ResponseWriter, r *http.Request) {
	// Load download folder from OS ENV
	destinationFolder := filepath.Clean("/mnt/downloads") // This is the docker directory, not the OS

	slog.Info("destinationFolder", "value", destinationFolder)
	// test the access to the directory
	if _, err := os.Stat(destinationFolder); os.IsNotExist(err) {
		slog.Warn("DESTINATION_FOLDER cannot be accessed")
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(rw, "Unable to load directory %s", destinationFolder)
		slog.Warn("Unable to load DESTINATION_FOLDER files")
		return
	}

	// load files here
	items, err := os.ReadDir(destinationFolder)
	if err != nil {
		slog.Warn("DESTINATION_FOLDER files cannot be accessed")
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(rw, "Unable to load directory %s", destinationFolder)
		return
	}
	// The json in which we marshal the slice of Files
	var Files []file
	// Create the json to return
	for _, item := range items {
		fileInfo, err := item.Info()
		if err != nil {
			slog.Warn("Error retrieving file information", "name", item.Name())
		}
		filetype := "file"
		if fileInfo.IsDir() {
			filetype = "dir"
		}
		f := file{
			Name:      item.Name(),
			Type:      filetype,
			Extension: ReturnExtension(item),
			Size:      FormatFileSize(fileInfo.Size()),
		}
		Files = append(Files, f)

	}
	// We now have Files, a slice of Files objects.
	// return the json list of files
	rw.Header().Set("Content-Type", "application/json")
	returnBody, err := json.Marshal(Files)
	if err != nil {
		slog.Warn(fmt.Sprintf("Error marshalling json: %s", err))
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(rw, err)
	}
	rw.Write(returnBody)
}

func DownloadFiles(rw http.ResponseWriter, r *http.Request) {
	// Load download folder from OS ENV
	destinationFolder := filepath.Clean(os.Getenv("DESTINATION_FOLDER"))
	// test the access to the directory
	if _, err := os.Stat(destinationFolder); os.IsNotExist(err) {
		slog.Warn("DESTINATION_FOLDER cannot be accessed")
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(rw, "Unable to load directory %s", destinationFolder)
		slog.Warn("Unable to load DESTINATION_FOLDER files")
		return
	}
	// This segment handles the downloads of files depending on the http request received. If not, serve the web UI

	// parse the request
	err := r.ParseForm()
	if err != nil {
		slog.Warn("Error parsing request:")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	// Ensure we received a value for file download
	requestedFileName := r.FormValue("filename")
	if requestedFileName == "" {
		slog.Info("Empty file name requested, cannot proceed")
		rw.WriteHeader(http.StatusNotFound)
		return
	} else {
		slog.Info("Received request to download", "file", requestedFileName, "Remote host", r.RemoteAddr)
	}
	// Security: ensure only files in the downloads are allowed
	cleanPath := filepath.Clean(requestedFileName)
	fullpath := filepath.Join(destinationFolder, cleanPath)

	if !strings.HasPrefix(fullpath, destinationFolder) {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte("Unable to download files outside of download directory: " + fullpath))
		return
	}
	// stream the file to the client
	rw.Header().Set("Content-Disposition", "attachment; filename=\""+filepath.Base(fullpath)+"\"")
	http.ServeFile(rw, r, fullpath)
}
