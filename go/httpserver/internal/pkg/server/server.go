package httpserver

import (
	"httpserver-cncamp/internal/pkg/server/filters"
	"net/http"

	log "k8s.io/klog/v2"
)

// ServerConfig represents the configuration for the server
type ServerConfig struct {
	Address     string
	Middlewares []filters.MiddlewareFunc //Middleware are executed in the order that they are applied to the Router.
	Router      http.Handler
}

type HttpServer struct {
	Server *http.Server
}

func NewServer(conf *ServerConfig) (*HttpServer, error) {
	configuredRouter := buildHandlerChain(conf)
	return &HttpServer{
		Server: &http.Server{
			Handler: configuredRouter,
			Addr:    conf.Address,
		}}, nil
}

// Run the Server
func (srv *HttpServer) Run() error {
	log.V(2).Info("Starting http server", srv.Server.Addr)
	return srv.Server.ListenAndServe()
}

// buildHandlerChain wraps the middleware with router
func buildHandlerChain(conf *ServerConfig) http.Handler {
	configuredRouter := conf.Router
	for _, adapter := range conf.Middlewares {
		configuredRouter = adapter(configuredRouter)
	}
	return configuredRouter
}
