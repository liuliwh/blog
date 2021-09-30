package main

import (
	"flag"
	"io"
	"net/http"
	_ "net/http/pprof"
	"strings"

	"github.com/golang/glog"
)

func main() {
	// set default flag value, the cmdline has higher priority
	flag.Set("v", "4")
	flag.Set("logtostderr", "true")
	flag.Parse()

	glog.V(2).Info("Starting http server...")

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/healthz", healthz)
	glog.Fatal(http.ListenAndServe(":8080", nil))
}

func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "200\n")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 接收客户端 request，并将 request 中带的 header 写入 response header
	// Echo the request headers to response header,
	// but exclude the representation headers
	for name, headers := range r.Header {
		k := http.CanonicalHeaderKey(name)
		// If request is POST, usually the request headers
		// contain Content-Type and Content-Length, which shouldn't
		// reprensent the response body content, so exlude those headers
		if !strings.HasPrefix(k, "Content-") {
			for _, h := range headers {
				w.Header().Add(k, h)
			}
		}
	}
	io.WriteString(w, "ok\n")
}
