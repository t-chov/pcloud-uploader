package pcloud

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// TestListFolder_success verifies that a valid response returns the expected metadata.
func TestListFolder_success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"result":0,"metadata":{"path":"/Backup","name":"Backup","isfolder":true,"contents":[]}}`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL)
	result, err := c.ListFolder("token", "/Backup", ListFolderOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Metadata.Name != "Backup" {
		t.Errorf("got name %q, want %q", result.Metadata.Name, "Backup")
	}
}

// TestListFolder_pathNormalization verifies that various path forms are normalized to /Backup.
func TestListFolder_pathNormalization(t *testing.T) {
	paths := []string{"Backup/", "/Backup", "Backup"}
	for _, p := range paths {
		t.Run(p, func(t *testing.T) {
			var gotPath string
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Query().Get("path")
				w.Write([]byte(`{"result":0,"metadata":{}}`))
			}))
			defer ts.Close()

			c := NewClient(ts.URL)
			_, _ = c.ListFolder("token", p, ListFolderOptions{})
			if gotPath != "/Backup" {
				t.Errorf("path %q: got %q, want %q", p, gotPath, "/Backup")
			}
		})
	}
}

// TestListFolder_queryParams verifies that auth, path, and option flags are sent correctly.
func TestListFolder_queryParams(t *testing.T) {
	var gotQuery url.Values
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		w.Write([]byte(`{"result":0,"metadata":{}}`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL)
	opts := ListFolderOptions{Recursive: true, NoFiles: true}
	_, _ = c.ListFolder("mytoken", "/Backup", opts)

	checks := map[string]string{
		"auth":      "mytoken",
		"path":      "/Backup",
		"recursive": "1",
		"nofiles":   "1",
	}
	for key, want := range checks {
		if got := gotQuery.Get(key); got != want {
			t.Errorf("query param %q: got %q, want %q", key, got, want)
		}
	}
}

// TestListFolder_httpError verifies that an HTTP failure returns an error.
func TestListFolder_httpError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	ts.Close()

	c := NewClient(ts.URL)
	_, err := c.ListFolder("token", "/Backup", ListFolderOptions{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
