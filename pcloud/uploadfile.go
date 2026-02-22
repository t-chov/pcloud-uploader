package pcloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// UploadFileResult represents the API response from uploadfile.
type UploadFileResult struct {
	Result    int            `json:"result"`
	Fileids   []int          `json:"fileids"`
	Checksums []Checksum     `json:"checksums"`
	Metadata  []FileMetadata `json:"metadata"`
}

// FileMetadata represents metadata for an uploaded file.
type FileMetadata struct {
	Ismine         bool    `json:"ismine"`
	Id             string  `json:"id"`
	Created        string  `json:"created"`
	IsShared       bool    `json:"isshared"`
	IsFolder       bool    `json:"isfolder"`
	Category       int     `json:"category"`
	ParentFolderId int     `json:"parentfolderid"`
	Icon           string  `json:"icon"`
	FileId         int     `json:"fileid"`
	Height         int     `json:"height"`
	Width          int     `json:"width"`
	Path           string  `json:"path"`
	Name           string  `json:"name"`
	ContentType    string  `json:"contenttype"`
	Size           int     `json:"size"`
	Thumb          bool    `json:"thumb"`
	Hash           big.Int `json:"hash"`
}

// Checksum holds file checksums returned by the API.
type Checksum struct {
	Sha1 string `json:"sha1"`
	Md5  string `json:"md5"`
}

// UploadFile uploads a file to the specified path on pCloud.
func (c *Client) UploadFile(auth, path string, file *os.File) (*UploadFileResult, error) {
	path = strings.TrimLeft(path, "/")
	path = "/" + path

	params := url.Values{
		"auth": {auth},
		"path": {filepath.Dir(path)},
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("filename", file.Name())
	if err != nil {
		return nil, fmt.Errorf("CreateFormFile: %v", err)
	}
	io.Copy(part, file)
	writer.Close()

	reqURL := fmt.Sprintf("%s/uploadfile?%s", c.BaseURL, params.Encode())
	req, err := http.NewRequest("POST", reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("newRequest: %v", err)
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("uploadfile: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read uploadfile body %s:%v", respBody, err)
	}

	var result UploadFileResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal uploadfile result %s:%v", string(respBody), err)
	}

	return &result, nil
}
