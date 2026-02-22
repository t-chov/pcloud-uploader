package pcloud

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestUploadFile_success verifies that a valid response returns the expected result.
func TestUploadFile_success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"result":0,"fileids":[1],"checksums":[{"sha1":"abc","md5":"def"}],"metadata":[]}`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL)
	result, err := c.UploadFile("token", "/Backup/foo.txt", "foo.txt", strings.NewReader("file content"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Result != 0 {
		t.Errorf("got result %d, want 0", result.Result)
	}
}

// TestUploadFile_pathExtraction verifies that the directory portion of path is sent as the path param.
func TestUploadFile_pathExtraction(t *testing.T) {
	var gotPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Query().Get("path")
		w.Write([]byte(`{"result":0,"metadata":[]}`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL)
	_, _ = c.UploadFile("token", "/Backup/foo.txt", "foo.txt", strings.NewReader(""))
	if gotPath != "/Backup" {
		t.Errorf("got path %q, want %q", gotPath, "/Backup")
	}
}

// TestUploadFile_httpError verifies that an HTTP failure returns an error.
func TestUploadFile_httpError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	ts.Close()

	c := NewClient(ts.URL)
	_, err := c.UploadFile("token", "/Backup/foo.txt", "foo.txt", strings.NewReader(""))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
