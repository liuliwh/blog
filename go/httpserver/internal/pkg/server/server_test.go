package httpserver

import (
	"httpserver-cncamp/internal/pkg/assert"
	"httpserver-cncamp/internal/pkg/server/filters"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	HelloClient string = "Hello, client"
	HeaderName1 string = "name1"
	HeaderName2 string = "name2"
	HeaderVal1  string = "Value1"
	HeaderVal2  string = "Value2"
)

func dummyHandler(message string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, message)
	})
}

func addHeader(name, value string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Add(name, value)
			h.ServeHTTP(rw, r)
		})
	}
}

func TestNewServer(t *testing.T) {
	router := dummyHandler(HelloClient)
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
		assert.StringEqual(t, w.Body.String(), HelloClient)
		headers := w.Result().Header

		assert.StringEqual(t, headers.Get(HeaderName1), HeaderVal1)
		assert.StringEqual(t, headers.Get(HeaderName2), HeaderVal2)
	})
	t.Run("NewServer should return the address specified", func(t *testing.T) {
		eIP, ePort, _ := net.SplitHostPort(ln.Addr().String())
		aIP, aPort, _ := net.SplitHostPort(s.Server.Addr)
		assert.StringEqual(t, aIP, eIP)
		assert.StringEqual(t, aPort, ePort)
	})
}

func TestAdapterHandler(t *testing.T) {
	const (
		HeaderName1 string = "name1"
		HeaderName2 string = "name2"
		HeaderVal1  string = "Value1"
		HeaderVal2  string = "Value2"
	)
	t.Run("NewServer should addheader1 & addheader2 to router", func(t *testing.T) {
		m1 := addHeader(HeaderName1, HeaderVal1)
		m2 := addHeader(HeaderName2, HeaderVal2)
		h := AdaptHandler(dummyHandler(HelloClient), m1, m2)
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		assert.StringEqual(t, w.Body.String(), HelloClient)
		headers := w.Result().Header

		assert.StringEqual(t, headers.Get(HeaderName1), HeaderVal1)
		assert.StringEqual(t, headers.Get(HeaderName2), HeaderVal2)
	})
}
