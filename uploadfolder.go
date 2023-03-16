package main

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

type uploadfileResult struct {
	Result    int                  `json:"result"`
	Fileids   []int                `json:"fileids"`
	Checksums []checksum           `json:"checksums"`
	Metadata  []uploadfileMetadata `json:"metadata"`
}

type uploadfileMetadata struct {
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

type checksum struct {
	Sha1 string `json:"sha1"`
	Md5  string `json:"md5"`
}

func uploadfile(auth, path string, file *os.File) error {
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
		return fmt.Errorf("CreateFormFile: %v", err)
	}
	io.Copy(part, file)
	defer writer.Close()

	url := fmt.Sprintf("https://api.pcloud.com/uploadfile?%s", params.Encode())
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("newRequest: %v", err)
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("uploadfile: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	// debug
	// fmt.Println(string(respBody))
	if err != nil {
		return fmt.Errorf("read uploadfile body %s:%v", respBody, err)
	}

	var result uploadfileResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("unmarshal uploadfile result %s:%v", string(respBody), err)
	}

	for _, metadata := range result.Metadata {
		fmt.Println(metadata.Hash.String())
	}

	return nil
}
