package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	appName = "pcloud-uploader"
	version = "0.0.1"
)

func main() {
	os.Exit(run())
}

func run() int {
	app := cli.NewApp()
	app.Name = appName
	app.Version = version
	app.Action = func(c *cli.Context) error {
		fmt.Println("new app!")
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		return 1
	} else {
		return 0
	}
}
