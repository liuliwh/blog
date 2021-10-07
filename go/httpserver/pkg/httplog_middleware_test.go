package httpserver

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	log "k8s.io/klog/v2"
)

func TestLoggedStatus(t *testing.T) {
	buf := setupLogOutput()
	h := setupHandler()
	s := httptest.NewServer(h)
	defer s.Close()
	client := s.Client()
	t.Cleanup(func() {
		log.Flush()
		buf.Reset()
		t.Log("clean")
	})
	executeGet := func(path string, wants ...string) {
		log.Flush()
		buf.Reset()
		client.Get(fmt.Sprint(s.URL, path))
		log.Flush()
		logStr := buf.String()
		for _, want := range wants {
			assertStringContains(t, logStr, want)
		}
	}

	t.Run("Should record status,ip,uri", func(t *testing.T) {
		executeGet("/", `srcIP="127.0.0.1"`, `URI="/`, `status=200`)
	})
	t.Run("Should record status,ip,uri after panic recovered", func(t *testing.T) {
		executeGet("/panic", `srcIP="127.0.0.1"`, `URI="/panic`, `status=500`)
	})
}

func assertStringContains(t *testing.T, s string, substr string) {
	if !strings.Contains(s, substr) {
		t.Errorf("want = %q, got = %q", substr, s)
	}
}

func setupLogOutput() *bytes.Buffer {
	log.InitFlags(nil)
	flag.Set("logtostderr", "false")
	flag.Set("alsologtostderr", "false")
	flag.Parse()

	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	return buf
}

func setupHandler() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})
	router.HandleFunc("/404", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
	})
	router.HandleFunc("/panic", func(rw http.ResponseWriter, r *http.Request) {
		panic("Test enter panic handler")
	})
	h := WithHttpLogging()(router)
	return h
}
