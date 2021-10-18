package filters

import (
	"fmt"
	"net/http"
)

// MiddlewareFunc is a function which receives an http.Handler and returns another http.Handler.
// Typically, the returned handler is a closure which does something with the http.ResponseWriter and http.Request passed
// to it, and then calls the handler passed as parameter to the MiddlewareFunc.
type MiddlewareFunc func(http.Handler) http.Handler

// EchoReqHeadersToResHeader add the request headers to response headers
// All the header name will be prefixed with X-Req- to avoid conflicting with
// response headers(especially the representation headers)
func EchoReqHeadersToResHeader() MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			for name, headers := range r.Header {
				for _, header := range headers {
					rw.Header().Add(fmt.Sprintf("X-Req-%s", name), header)
				}
			}
			h.ServeHTTP(rw, r)
		})
	}
}

// AddRespHeader is the middlware to add the name/value to the response header
func AddRespHeader(name, value string) MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Add(name, value)
			h.ServeHTTP(rw, r)
		})
	}
}
