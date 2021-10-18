package appserver

import (
	"io"
	"net/http"
)

// healthz is the handler for application level health check.
func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "200\n")
}

// rootHandler is the handler for http server
func rootHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok\n")
}
