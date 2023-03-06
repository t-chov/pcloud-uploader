package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type auth struct {
	Auth string `json:"auth"`
}

func getAuth(username, password string) (*string, error) {
	params := []string{
		fmt.Sprintf("username=%s", username),
		fmt.Sprintf("password=%s", password),
	}
	url := fmt.Sprintf("https://api.pcloud.com/userinfo?getauth=1&logout=1&%s", strings.Join(params, "&"))
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("getauth: %v", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("load auth body: %v", err)
	}
	var auth auth
	if err := json.Unmarshal(body, &auth); err != nil {
		return nil, fmt.Errorf("parse auth `%s`: %v", body, err)
	}
	return &auth.Auth, nil
}
