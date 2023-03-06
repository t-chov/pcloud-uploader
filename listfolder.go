package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func listfolder(auth, path string, recursive bool) error {
	var recursiveInt int
	if recursive {
		recursiveInt = 1
	}
	params := []string{
		fmt.Sprintf("auth=%s", auth),
		fmt.Sprintf("path=%s", path),
		fmt.Sprintf("recursive=%d", recursiveInt),
	}
	url := fmt.Sprintf("https://api.pcloud.com/listfolder?%s", strings.Join(params, "&"))
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("get listfolder: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read listfolder body: %v", err)
	}
	fmt.Println(string(body))
	return nil
}
