package appserver

import (
	httpserver "httpserver-cncamp/internal/pkg/server"
	"httpserver-cncamp/internal/pkg/server/filters"
	"net/http"
	"os"
)

// NewServer setup the server config with the given listen address
func NewServer(address string) (*httpserver.HttpServer, error) {
	c := &httpserver.ServerConfig{
		Address:     address,
		Router:      setupRoutes(),
		Middlewares: setupMiddlewares(),
	}
	return httpserver.NewServer(c)
}

// setupRoutes use DefaultServeMux to utilize the pprof sideeffect
func setupRoutes() *http.ServeMux {
	router := http.DefaultServeMux
	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/healthz", healthz)
	return router
}

func setupMiddlewares() []filters.MiddlewareFunc {
	version := os.Getenv("VERSION")
	res := make([]filters.MiddlewareFunc, 0)
	res = append(res, filters.EchoReqHeadersToResHeader())
	res = append(res, filters.AddRespHeader("version", version))
	res = append(res, filters.WithHttpLogging())
	return res
}
