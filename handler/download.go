package handler

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const destinationFolder = "/mnt/downloads"

type FSEntry struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	IsDir     bool      `json:"isDir"`
	Size      string    `json:"size,omitempty"`
	Extension string    `json:"extension,omitempty"`
	Children  []FSEntry `json:"children,omitempty"`
}

func ReturnExtension(s fs.DirEntry) string {
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
	slog.Info("File listing request", "remote", r.RemoteAddr, "folder", destinationFolder)
	// test the access to the directory
	if _, err := os.Stat(destinationFolder); os.IsNotExist(err) {
		slog.Warn("downloads folder cannot be accessed")
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Unable to load directory %s", destinationFolder)
		slog.Warn("Unable to load downloads folder files")
		return
	}
	// Create a map to store folders and their files in it.
	folders := make(map[string][]FSEntry)
	// load files using walkPath to get all underlying elements
	err := filepath.WalkDir(destinationFolder, func(path string, d fs.DirEntry, err error) error {
		if path == destinationFolder {
			// this is  the root folder, skip it
			return nil
		}
		if d.IsDir() {
			depth := strings.Count(strings.TrimPrefix(path, destinationFolder), string(filepath.Separator))
			if depth > 1 {
				return fs.SkipDir // do not go below 1 level in the tree.
			}
			folders[d.Name()] = []FSEntry{}
			return nil
		}
		parentFolder := filepath.Base(filepath.Dir(path))

		fileInfo, err := d.Info()
		if err != nil {
			slog.Warn(err.Error())
		}
		relativePath, _ := filepath.Rel(destinationFolder, path)

		folders[parentFolder] = append(folders[parentFolder], FSEntry{
			Name:      d.Name(),
			Path:      relativePath,
			IsDir:     d.IsDir(),
			Size:      FormatFileSize(fileInfo.Size()),
			Extension: ReturnExtension(d),
		})
		return nil
	})
	if err != nil {
		slog.Warn("Error retrieving directory files", "err", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// We now have folders, a map of folders and their sub items.
	// return the json list of files
	rw.Header().Set("Content-Type", "application/json")
	returnBody, err := json.Marshal(folders)
	if err != nil {
		slog.Warn(fmt.Sprintf("Error marshalling json: %s", err))
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(rw, err)
	}
	rw.Write(returnBody)
}

func DownloadFiles(rw http.ResponseWriter, r *http.Request) {
	// test the access to the directory
	if _, err := os.Stat(destinationFolder); os.IsNotExist(err) {
		slog.Warn("downloads folder cannot be accessed")
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(rw, "Unable to load directory %s", destinationFolder)
		slog.Warn("Unable to load downloads folder files")
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

func DownloadFolder(rw http.ResponseWriter, r *http.Request) {
	// parse the request
	err := r.ParseForm()
	if err != nil {
		slog.Warn("Error parsing request:")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	archiveName := r.FormValue("folder")
	slog.Info("Received request to download folder", "name", archiveName, "remote", r.RemoteAddr)

	// build the path to the folder
	FolderPath := filepath.Join(destinationFolder, archiveName)
	// test the access to the directory
	if _, err := os.Stat(FolderPath); os.IsNotExist(err) {
		slog.Warn(FolderPath + "cannot be accessed")
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(rw, "Unable to load directory %s", FolderPath)
		return
	}
	// prevent path traversal
	if !strings.HasPrefix(FolderPath, destinationFolder) {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte("Access denied"))
		return
	}
	rw.Header().Set("Content-Type", "application/zip")
	rw.Header().Set("Content-Disposition",
		"attachment; filename=\""+archiveName+".zip\"")
	// Create the archive
	zipWriter := zip.NewWriter(rw)
	defer zipWriter.Close()
	// Create a buffer for io.reader
	buffer := make([]byte, 4*1024*1024) // 4 MB
	filepath.WalkDir(FolderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			slog.Warn("WalkDir error", "path", path, "err", err)
			return err
		}
		if d.IsDir() {
			return nil // les dossiers sont créés automatiquement par le zip
		}

		// relative path, which is name and structure in archives
		relativePath, err := filepath.Rel(FolderPath, path)
		if err != nil {
			slog.Warn(err.Error())
		}
		// create the entry
		zipEntry, err := zipWriter.CreateHeader(&zip.FileHeader{
			Name:   relativePath,
			Method: zip.Store,
		})
		if err != nil {
			slog.Warn(err.Error())
		}
		// copy the file
		f, err := os.Open(path)
		if err != nil {
			slog.Warn(err.Error())
		}
		defer f.Close()
		io.CopyBuffer(zipEntry, f, buffer)

		return nil
	})
}
