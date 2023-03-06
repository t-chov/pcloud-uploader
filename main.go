package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

const (
	APP_NAME            = "pcloud"
	VERSION             = "0.0.1"
	ENV_PCLOUD_USERNAME = "PCLOUD_USERNAME"
	ENV_PCLOUD_PASSWORD = "PCLOUD_PASSWORD"
)

func main() {
	os.Exit(run())
}

func run() int {
	app := cli.NewApp()
	app.Name = APP_NAME
	app.Version = VERSION
	app.Action = appAction
	if err := app.Run(os.Args); err != nil {
		return 1
	} else {
		return 0
	}
}

func appAction(c *cli.Context) error {
	var username, password, auth *string
	var err error
	if username, err = loadFromInput(ENV_PCLOUD_USERNAME, "username"); err != nil {
		return err
	}
	if password, err = loadFromInput(ENV_PCLOUD_PASSWORD, "password"); err != nil {
		return err
	}

	if auth, err = getAuth(*username, *password); err != nil {
		return err
	} else if *auth == "" {
		return fmt.Errorf("empty auth")
	}
	fmt.Printf("auth: %s\n", *auth)
	return nil
}

func loadFromInput(envName string, title string) (*string, error) {
	value := os.Getenv(envName)
	if value == "" {
		fmt.Printf("%s: ", title)
		if _, err := fmt.Scan(&value); err != nil {
			return nil, fmt.Errorf("scan %s: %v", title, err)
		}
	}
	value = strings.TrimSpace(value)
	return &value, nil
}
