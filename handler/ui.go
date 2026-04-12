package handler

import (
	_ "embed"
	"net/http"
)

//go:embed static/index.html
var indexHTML []byte

func HandleUI(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	rw.Write(indexHTML)
}
