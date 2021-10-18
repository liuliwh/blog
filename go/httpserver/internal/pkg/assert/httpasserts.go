package assert

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func ResponseCode(t *testing.T, err error, resp *http.Response, statusCodeExpected int) {
	ErrIsNil(t, err)
	if resp.StatusCode != statusCodeExpected {
		t.Errorf("Expected Status Code: %d, but got: %d", statusCodeExpected, resp.StatusCode)
	}
}

func ResponseBody(t *testing.T, err error, resp *http.Response, expectedBody string) {
	ErrIsNil(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Errorf("unexpected error reading body: %v", err)
	}
	if !bytes.Equal(body, []byte(expectedBody)) {
		t.Errorf("response should be %q, was: %q", expectedBody, string(body))
	}
}
