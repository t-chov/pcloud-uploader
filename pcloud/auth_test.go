package pcloud

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAuthenticate_success verifies that a valid JSON response returns the token.
func TestAuthenticate_success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"auth":"token123"}`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL)
	token, err := c.Authenticate("user", "pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "token123" {
		t.Errorf("got token %q, want %q", token, "token123")
	}
}

// TestAuthenticate_httpError verifies that an HTTP failure returns an error.
func TestAuthenticate_httpError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	ts.Close() // close immediately so connections are refused

	c := NewClient(ts.URL)
	_, err := c.Authenticate("user", "pass")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestAuthenticate_invalidJSON verifies that malformed JSON returns a parse error.
func TestAuthenticate_invalidJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`not-json`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL)
	_, err := c.Authenticate("user", "pass")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
