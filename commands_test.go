package main

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/t-chov/pcloud-uploader/pcloud"
)

type mockAPI struct {
	authenticateFn func(u, p string) (string, error)
	listFolderFn   func(a, path string, opts pcloud.ListFolderOptions) (*pcloud.ListFolderResult, error)
	uploadFileFn   func(a, path, filename string, r io.Reader) (*pcloud.UploadFileResult, error)
}

func (m *mockAPI) Authenticate(u, p string) (string, error) {
	return m.authenticateFn(u, p)
}

func (m *mockAPI) ListFolder(a, path string, opts pcloud.ListFolderOptions) (*pcloud.ListFolderResult, error) {
	return m.listFolderFn(a, path, opts)
}

func (m *mockAPI) UploadFile(a, path, filename string, r io.Reader) (*pcloud.UploadFileResult, error) {
	return m.uploadFileFn(a, path, filename, r)
}

// TestRunWithClient_ls verifies that the ls command succeeds and calls ListFolder.
func TestRunWithClient_ls(t *testing.T) {
	t.Setenv(ENV_PCLOUD_USERNAME, "user")
	t.Setenv(ENV_PCLOUD_PASSWORD, "pass")

	called := false
	mock := &mockAPI{
		authenticateFn: func(u, p string) (string, error) {
			return "testtoken", nil
		},
		listFolderFn: func(a, path string, opts pcloud.ListFolderOptions) (*pcloud.ListFolderResult, error) {
			called = true
			return &pcloud.ListFolderResult{
				Metadata: pcloud.FolderMetadata{
					Path:     "/Backup",
					Name:     "Backup",
					IsFolder: true,
				},
			}, nil
		},
	}

	code := runWithClient(mock, []string{"app", "ls", "/Backup"})
	if code != 0 {
		t.Errorf("got exit code %d, want 0", code)
	}
	if !called {
		t.Error("ListFolder was not called")
	}
}

// TestRunWithClient_up verifies that the up command succeeds and calls UploadFile.
func TestRunWithClient_up(t *testing.T) {
	t.Setenv(ENV_PCLOUD_USERNAME, "user")
	t.Setenv(ENV_PCLOUD_PASSWORD, "pass")

	f, err := os.CreateTemp(t.TempDir(), "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("hello")
	f.Close()

	called := false
	mock := &mockAPI{
		authenticateFn: func(u, p string) (string, error) {
			return "testtoken", nil
		},
		uploadFileFn: func(a, path, filename string, r io.Reader) (*pcloud.UploadFileResult, error) {
			called = true
			return &pcloud.UploadFileResult{}, nil
		},
	}

	code := runWithClient(mock, []string{"app", "up", f.Name(), "/Backup"})
	if code != 0 {
		t.Errorf("got exit code %d, want 0", code)
	}
	if !called {
		t.Error("UploadFile was not called")
	}
}

// TestRunWithClient_authError verifies that an authentication failure returns exit code 1.
func TestRunWithClient_authError(t *testing.T) {
	t.Setenv(ENV_PCLOUD_USERNAME, "user")
	t.Setenv(ENV_PCLOUD_PASSWORD, "pass")

	mock := &mockAPI{
		authenticateFn: func(u, p string) (string, error) {
			return "", fmt.Errorf("auth failed")
		},
	}

	code := runWithClient(mock, []string{"app", "ls", "/Backup"})
	if code != 1 {
		t.Errorf("got exit code %d, want 1", code)
	}
}
