package handler

import (
	_ "embed"
	"net/http"
)

//go:embed static/index.html
var indexHTML []byte

func HandleUI(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write(indexHTML)
}
