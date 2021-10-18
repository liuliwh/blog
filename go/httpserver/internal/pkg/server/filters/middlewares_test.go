package filters

import (
	"fmt"
	"httpserver-cncamp/internal/pkg/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEchoReqHeadersToResHeader(t *testing.T) {
	const (
		HEADER string = "header"
		VAL    string = "val"
		PREFIX string = "X-Req-"
	)
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})
	h := EchoReqHeadersToResHeader()(router)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add(HEADER, VAL)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	actual := w.Result().Header.Get(fmt.Sprint(PREFIX, HEADER))
	assert.StringEqual(t, actual, VAL)
}
