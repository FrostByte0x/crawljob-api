package handler

import (
	_ "embed"
	"log/slog"
	"net/http"
)

//go:embed static/download.html
var downloadInterfacetype []byte

func HandleDownloadUi(rw http.ResponseWriter, request *http.Request) {
	slog.Info("Download UI request", "remote", request.RemoteAddr)
	rw.Header().Set("Content-Type", "text/html")
	rw.Write(downloadInterfacetype)
}
