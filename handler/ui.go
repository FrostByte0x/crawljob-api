package handler

import (
	_ "embed"
	"log/slog"
	"net/http"
)

//go:embed static/index.html
var indexHTML []byte

func HandleUI(rw http.ResponseWriter, request *http.Request) {
	slog.Info("UI request", "remote", request.RemoteAddr)
	rw.Header().Set("Content-Type", "text/html")
	rw.Write(indexHTML)
}
