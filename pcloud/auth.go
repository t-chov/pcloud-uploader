package pcloud

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type authResponse struct {
	Auth string `json:"auth"`
}

// Authenticate obtains a session token from the pCloud API.
func (c *Client) Authenticate(username, password string) (string, error) {
	params := []string{
		fmt.Sprintf("username=%s", username),
		fmt.Sprintf("password=%s", password),
	}
	url := fmt.Sprintf("%s/userinfo?getauth=1&logout=1&%s", c.BaseURL, strings.Join(params, "&"))
	response, err := c.HTTPClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("getauth: %v", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("load auth body: %v", err)
	}
	var auth authResponse
	if err := json.Unmarshal(body, &auth); err != nil {
		return "", fmt.Errorf("parse auth `%s`: %v", body, err)
	}
	return auth.Auth, nil
}
