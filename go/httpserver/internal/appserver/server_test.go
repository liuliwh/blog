package appserver

import (
	"fmt"
	"httpserver-cncamp/internal/pkg/assert"
	httpserver "httpserver-cncamp/internal/pkg/server"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewServerConf(t *testing.T) {
	t.Setenv("VERSION", "v1")
	s, srvHolder := setUpUnstartedServer(t)
	s.Start()
	defer s.Close()
	client := s.Client()
	const HEADER string = "Name"

	executeGet := func(path string) (*http.Response, error) {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprint(s.URL, path), nil)
		req.Header.Add(HEADER, HEADER)
		req.Header.Set("User-Agent", "UA")
		return client.Do(req)
	}

	t.Run("The address should set as specified", func(t *testing.T) {
		assert.StringEqual(t, srvHolder.Server.Addr, s.Config.Addr)
	})

	t.Run("root handler is set", func(t *testing.T) {
		w, err := executeGet("/")
		assert.ResponseBody(t, err, w, "ok\n")
	})
	t.Run("healthz handler is set", func(t *testing.T) {
		w, err := executeGet("/healthz")
		assert.ResponseCode(t, err, w, http.StatusOK)
	})
	t.Run("env VERSION should be added into resp header", func(t *testing.T) {
		w, _ := executeGet("/")
		assert.StringEqual(t, w.Header.Get("Version"), "v1")
	})
	t.Run("request headers should be added into resp header", func(t *testing.T) {
		w, _ := executeGet("/")
		assert.StringEqual(t, w.Header.Get(fmt.Sprint("X-Req-", HEADER)), HEADER)
		assert.StringEqual(t, w.Header.Get(fmt.Sprint("X-Req-", "User-Agent")), w.Request.UserAgent())
	})
}

func setUpUnstartedServer(t *testing.T) (*httptest.Server, *httpserver.HttpServer) {
	// Find the port
	ln, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		t.Log("Error on listening")
	}
	srvHolder, _ := NewServer(ln.Addr().String())

	s := httptest.NewUnstartedServer(srvHolder.Server.Handler)
	s.Config.Addr = ln.Addr().String()
	return s, srvHolder
}
