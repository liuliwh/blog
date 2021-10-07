package httpserver

import (
	"net"
	"net/http"
	"runtime"

	log "k8s.io/klog/v2"
)

// responseLogger is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseLogger struct {
	w              http.ResponseWriter
	req            *http.Request
	status         int
	statusRecorded bool
}

func (rl *responseLogger) Status() int {
	return rl.status
}

// WriteHeader implements http.ResponseWriter.
func (rl *responseLogger) WriteHeader(code int) {
	rl.recordStatus(code)
	rl.w.WriteHeader(code)
}

// Header implements http.ResponseWriter.
func (rl *responseLogger) Header() http.Header {
	return rl.w.Header()
}

// Write implements http.ResponseWriter.
func (rl *responseLogger) Write(b []byte) (int, error) {
	// based on "If WriteHeader has not yet been called, Write calls
	// WriteHeader(http.StatusOK) before writing the data."
	// so if status is not Recorded, then we need set default value
	// to our wrapper struct
	if !rl.statusRecorded {
		rl.recordStatus(http.StatusOK)
	}
	return rl.w.Write(b)
}
func (rl *responseLogger) recordStatus(status int) {
	rl.status = status
	rl.statusRecorded = true
}
func (rl *responseLogger) Log() {
	ip, _, _ := net.SplitHostPort(rl.req.RemoteAddr)
	keysAndValues := []interface{}{
		"srcIP", ip,
		"URI", rl.req.RequestURI,
	}
	if err := recover(); err != nil {
		rl.WriteHeader(http.StatusInternalServerError)
		rl.Write([]byte("Oops, something went wrong\n"))
		stack := make([]byte, 50*1024)
		stack = stack[:runtime.Stack(stack, false)]
		keysAndValues = append(keysAndValues, "stack", stack)
	}
	keysAndValues = append(keysAndValues, "status", rl.status)
	log.InfoS(
		"HTTP",
		keysAndValues...,
	)
}

// WithHttpLogging logs the incoming HTTP request & its duration.
func WithHttpLogging() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapped := &responseLogger{
				w:   w,
				req: r,
			}
			defer wrapped.Log()

			next.ServeHTTP(wrapped, r)

		})
	}
}
