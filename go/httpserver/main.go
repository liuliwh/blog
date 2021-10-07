package main

import (
	"flag"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"

	httpserver "github.com/liuliwh/blog/go/httpserver/pkg"
	log "k8s.io/klog/v2"
)

func main() {
	// Set up logging
	log.InitFlags(nil)

	// get config
	var address string
	flag.StringVar(&address, "address", ":8080", "HTTP Server Address")
	flag.Set("v", "4")
	flag.Set("logtostderr", "true")
	flag.Parse()

	c := newServerConf(address)
	srv, _ := httpserver.NewServer(c)
	err := srv.Run()
	log.Fatalf("Couldn't run: %s", err)
}

// newServerConf setup the server config with the given listen address
func newServerConf(address string) *httpserver.ServerConfig {
	return &httpserver.ServerConfig{
		Address:     address,
		Router:      setupRoutes(),
		Middlewares: setupMiddlewares(),
	}
}

// setupRoutes use DefaultServeMux to utilize the pprof sideeffect
func setupRoutes() *http.ServeMux {
	router := http.DefaultServeMux
	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/healthz", healthz)
	return router
}

func setupMiddlewares() []httpserver.MiddlewareFunc {
	version := os.Getenv("VERSION")
	res := make([]httpserver.MiddlewareFunc, 0)
	res = append(res, httpserver.EchoReqHeadersToResHeader())
	res = append(res, httpserver.AddRespHeader("version", version))
	res = append(res, httpserver.WithHttpLogging())
	return res
}

// healthz is the handler for application level health check.
func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "200\n")
}

// rootHandler is the handler for http server
func rootHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok\n")
}
