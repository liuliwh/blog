package appserver

import (
	"httpserver-cncamp/internal/pkg/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlersGet(t *testing.T) {
	executeGet := func(path, expected string, handler http.HandlerFunc) {
		req, _ := http.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()
		handler(w, req)
		assert.StringEqual(t, w.Body.String(), expected)
	}

	t.Run("root handler", func(t *testing.T) {
		executeGet("/", "ok\n", rootHandler)
	})
	t.Run("healthz handler", func(t *testing.T) {
		executeGet("/healthz", "200\n", healthz)
	})
}
