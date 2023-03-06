package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type listfolderResult struct {
	Result   int                `json:"result"`
	Metadata listfolderMetadata `json:"metadata"`
}

type listfolderMetadata struct {
	Path     string               `json:"path"`
	Name     string               `json:"name"`
	Created  string               `json:"created"`
	IsMine   bool                 `json:"ismine"`
	Thumb    bool                 `json:"thumb"`
	Modified string               `json:"Modified"`
	Id       string               `json:"id"`
	IsShared bool                 `json:"isshared"`
	Icon     string               `json:"icon"`
	IsFolder bool                 `json:"isfolder"`
	FolderId int                  `json:"folderid"`
	Contents []listfolderMetadata `json:"contents"`
}

var printFolder = color.New(color.FgBlue).PrintlnFunc()

func listfolder(auth, path string, recursive, showdeleted, nofiles, noshares bool) error {
	// must start with `/`, must remove end of `/`
	path = strings.Trim(path, "/")
	path = "/" + path

	params := url.Values{
		"auth":        {auth},
		"path":        {path},
		"recursive":   {strconv.Itoa(btoi(recursive))},
		"showdeleted": {strconv.Itoa(btoi(showdeleted))},
		"nofiles":     {strconv.Itoa(btoi(nofiles))},
		"noshares":    {strconv.Itoa(btoi(noshares))},
	}
	url := fmt.Sprintf("https://api.pcloud.com/listfolder?%s", params.Encode())

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("get listfolder: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read listfolder body: %v", err)
	}
	// debug
	// fmt.Println(string(body))

	var result listfolderResult
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("unmarshal listfolder result `%s`: %v", string(body), err)
	}
	printListfolder("/", result.Metadata)
	return nil
}

func printListfolder(path string, ls listfolderMetadata) {
	if ls.Path == "/" {
		printFolder("/")
	} else if ls.IsFolder {
		path += ls.Name + "/"
		printFolder(path)
	} else {
		path += ls.Name
		fmt.Println(path)
	}
	for _, content := range ls.Contents {
		printListfolder(path, content)
	}
}
