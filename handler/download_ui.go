package handler

import (
	_ "embed"
	"net/http"
)

//go:embed static/download.html
var downloadInterfacetype []byte

func HandleDownloadUi(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	rw.Write(downloadInterfacetype)
}
