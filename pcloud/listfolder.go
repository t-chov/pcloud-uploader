package pcloud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

// ListFolderResult represents the API response from listfolder.
type ListFolderResult struct {
	Result   int            `json:"result"`
	Metadata FolderMetadata `json:"metadata"`
}

// FolderMetadata represents a file or folder entry returned by the API.
type FolderMetadata struct {
	Path     string           `json:"path"`
	Name     string           `json:"name"`
	Created  string           `json:"created"`
	Ismine   bool             `json:"ismine"`
	Thumb    bool             `json:"thumb"`
	Modified string           `json:"Modified"`
	Id       string           `json:"id"`
	IsShared bool             `json:"isshared"`
	Icon     string           `json:"icon"`
	IsFolder bool             `json:"isfolder"`
	FolderId int              `json:"folderid"`
	Contents []FolderMetadata `json:"contents"`
}

// ListFolder retrieves the contents of a folder from the pCloud API.
func (c *Client) ListFolder(auth, path string, opts ListFolderOptions) (*ListFolderResult, error) {
	// must start with `/`, must remove end of `/`
	path = strings.Trim(path, "/")
	path = "/" + path

	params := url.Values{
		"auth":        {auth},
		"path":        {path},
		"recursive":   {strconv.Itoa(btoi(opts.Recursive))},
		"showdeleted": {strconv.Itoa(btoi(opts.ShowDeleted))},
		"nofiles":     {strconv.Itoa(btoi(opts.NoFiles))},
		"noshares":    {strconv.Itoa(btoi(opts.NoShares))},
	}
	reqURL := fmt.Sprintf("%s/listfolder?%s", c.BaseURL, params.Encode())

	resp, err := c.HTTPClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("get listfolder: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read listfolder body: %v", err)
	}

	var result ListFolderResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal listfolder result %s: %v", string(body), err)
	}
	return &result, nil
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
