package httpserver

import (
	"fmt"
	"httpserver-cncamp/internal/pkg/assert"
	"httpserver-cncamp/internal/pkg/server/filters"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewServer(t *testing.T) {
	const (
		HelloClient string = "Hello, client"
		HeaderName1 string = "name1"
		HeaderName2 string = "name2"
		HeaderVal1  string = "Value1"
		HeaderVal2  string = "Value2"
	)
	assertStringEqual := assert.StringEqual
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, HelloClient)
	})
	addHeader := func(name, value string) func(http.Handler) http.Handler {
		return func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.Header().Add(name, value)
				h.ServeHTTP(rw, r)
			})
		}
	}
	m1 := addHeader(HeaderName1, HeaderVal1)
	m2 := addHeader(HeaderName2, HeaderVal2)
	ln, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		t.Fatalf("failed to listen on %v: %v", "0.0.0.0:0", err)
	}
	c := ServerConfig{
		Address:     ln.Addr().String(),
		Router:      router,
		Middlewares: []filters.MiddlewareFunc{m1, m2},
	}
	s, _ := NewServer(&c)

	t.Run("NewServer should addheader1 & addheader2 to router", func(t *testing.T) {
		h := s.Server.Handler
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		assertStringEqual(t, w.Body.String(), HelloClient)
		headers := w.Result().Header

		assertStringEqual(t, headers.Get(HeaderName1), HeaderVal1)
		assertStringEqual(t, headers.Get(HeaderName2), HeaderVal2)
	})
	t.Run("NewServer should return the address specified", func(t *testing.T) {

		eIP, ePort, _ := net.SplitHostPort(ln.Addr().String())
		aIP, aPort, _ := net.SplitHostPort(s.Server.Addr)
		assertStringEqual(t, aIP, eIP)
		assertStringEqual(t, aPort, ePort)
	})
}
