package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"

	log "k8s.io/klog/v2"
)

// MiddlewareFunc is a function which receives an http.Handler and returns another http.Handler.
// Typically, the returned handler is a closure which does something with the http.ResponseWriter and http.Request passed
// to it, and then calls the handler passed as parameter to the MiddlewareFunc.
type MiddlewareFunc func(http.Handler) http.Handler

// echoReqHeadersToResHeader add the request headers to response headers
func echoReqHeadersToResHeader() MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			log.Info("echoReqHeadersToResHeader")
			for name, headers := range r.Header {
				for _, header := range headers {
					rw.Header().Add(fmt.Sprintf("X-Req-%s", name), header)
				}
			}
			h.ServeHTTP(rw, r)
		})
	}
}

func addHeader(name, value string) MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			log.Info("addHeader")
			rw.Header().Add(name, value)
			h.ServeHTTP(rw, r)
		})
	}
}

func main() {
	// Set up logging
	log.InitFlags(nil)

	// get config
	var address string
	flag.StringVar(&address, "address", ":8080", "HTTP Server Address")
	flag.Set("v", "4")
	flag.Set("logtostderr", "true")
	flag.Parse()

	// Set up the routes and middleware
	router := setupRoutes()
	version := os.Getenv("VERSION")

	runServer(serverConfig{
		Address:     address,
		Router:      router,
		Middlewares: []MiddlewareFunc{echoReqHeadersToResHeader(), addHeader("version", version)},
	})
}

// serverConfig represents the configuration for the server
type serverConfig struct {
	Address     string
	Middlewares []MiddlewareFunc //Middleware are executed in the order that they are applied to the Router.
	Router      http.Handler
}

// setupRoutes use DefaultServeMux to utilize the pprof sideeffect
func setupRoutes() *http.ServeMux {
	router := http.DefaultServeMux
	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/healthz", healthz)
	return router
}
func runServer(conf serverConfig) {
	srv, err := newServer(conf)
	if err != nil {
		log.Fatal(err)
	}
	log.V(2).Info("Starting http server", conf.Address)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error(err)
	}
}

// newServer applies the middlewares to the router, and return the configured
// handler to http.server
func newServer(conf serverConfig) (*http.Server, error) {
	// config the middleware to the router
	configuredRouter := conf.Router
	for _, adapter := range conf.Middlewares {
		configuredRouter = adapter(configuredRouter)
	}
	return &http.Server{
		Handler: configuredRouter,
		Addr:    conf.Address,
	}, nil
}

// healthz is the handler for application level health check.
func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "200\n")
}

// rootHandler is the handler for http server
func rootHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok\n")
}
