package pcloud

import (
	"net/http"
	"os"
)

const defaultBaseURL = "https://api.pcloud.com"

// API defines the pCloud API operations.
type API interface {
	Authenticate(username, password string) (string, error)
	ListFolder(auth, path string, opts ListFolderOptions) (*ListFolderResult, error)
	UploadFile(auth, path string, file *os.File) (*UploadFileResult, error)
}

// Client implements the API interface using HTTP calls to the pCloud API.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new pCloud API client. If baseURL is empty, the default
// pCloud API URL is used.
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

// ListFolderOptions controls the behavior of a ListFolder call.
type ListFolderOptions struct {
	Recursive   bool
	ShowDeleted bool
	NoFiles     bool
	NoShares    bool
}
