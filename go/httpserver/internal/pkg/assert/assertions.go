package assert

import "testing"

func StringEqual(t testing.TB, actual, expected string) {
	t.Helper()
	if actual != expected {
		t.Errorf("response body is wrong, got %q want %q", actual, expected)
	}
}
func ErrIsNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}
}
