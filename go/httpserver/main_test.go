package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	httpserver "github.com/liuliwh/blog/go/httpserver/pkg"
)

func TestHandlersGet(t *testing.T) {
	executeGet := func(path, expected string, handler http.HandlerFunc) {
		req, _ := http.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()
		handler(w, req)
		stringEqual(t, w.Body.String(), expected)
	}

	t.Run("root handler", func(t *testing.T) {
		executeGet("/", "ok\n", rootHandler)
	})
	t.Run("healthz handler", func(t *testing.T) {
		executeGet("/healthz", "200\n", healthz)
	})
}

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
		stringEqual(t, srvHolder.Server.Addr, s.Config.Addr)
	})

	t.Run("root handler is set", func(t *testing.T) {
		w, err := executeGet("/")
		assertResponseBody(t, err, w, "ok\n")
	})
	t.Run("healthz handler is set", func(t *testing.T) {
		w, err := executeGet("/healthz")
		assertResponseCode(t, err, w, http.StatusOK)
	})
	t.Run("env VERSION should be added into resp header", func(t *testing.T) {
		w, _ := executeGet("/")
		stringEqual(t, w.Header.Get("Version"), "v1")
	})
	t.Run("request headers should be added into resp header", func(t *testing.T) {
		w, _ := executeGet("/")
		stringEqual(t, w.Header.Get(fmt.Sprint("X-Req-", HEADER)), HEADER)
		stringEqual(t, w.Header.Get(fmt.Sprint("X-Req-", "User-Agent")), w.Request.UserAgent())
	})
}

func assertResponseCode(t *testing.T, err error, resp *http.Response, statusCodeExpected int) {
	assertErrIsNil(t, err)
	if resp.StatusCode != statusCodeExpected {
		t.Errorf("Expected Status Code: %d, but got: %d", statusCodeExpected, resp.StatusCode)
	}
}

func assertResponseBody(t *testing.T, err error, resp *http.Response, expectedBody string) {
	assertErrIsNil(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Errorf("unexpected error reading body: %v", err)
	}
	if !bytes.Equal(body, []byte(expectedBody)) {
		t.Errorf("response should be %q, was: %q", expectedBody, string(body))
	}
}

func assertErrIsNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}
}
func stringEqual(t testing.TB, actual, expected string) {
	if actual != expected {
		t.Errorf("response body is wrong, got %q want %q", actual, expected)
	}
}
func setUpUnstartedServer(t *testing.T) (*httptest.Server, *httpserver.HttpServer) {
	// Find the port
	ln, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		t.Log("Error on listening")
	}
	t.Logf("addr=%v", ln.Addr().String())
	conf := newServerConf(ln.Addr().String())
	srvHolder, _ := httpserver.NewServer(conf)

	s := httptest.NewUnstartedServer(srvHolder.Server.Handler)
	s.Config.Addr = ln.Addr().String()
	return s, srvHolder
}
