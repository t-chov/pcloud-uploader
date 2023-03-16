package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

const (
	APP_NAME            = "pcloud"
	VERSION             = "0.0.1"
	ENV_PCLOUD_USERNAME = "PCLOUD_USERNAME"
	ENV_PCLOUD_PASSWORD = "PCLOUD_PASSWORD"
)

var commands = []*cli.Command{
	{
		Name:    "ls",
		Aliases: []string{"listfolder"},
		Usage:   "Receive data for a folder.",
		Action: func(c *cli.Context) error {
			auth := c.Context.Value("auth").(string)
			if err := listfolder(auth, c.Args().First(), c.Bool("recursive"), c.Bool("showdeleted"), c.Bool("nofiles"), c.Bool("noshares")); err != nil {
				return printError(err)
			}
			return nil
		},
		ArgsUsage: "<PATH>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "recursive",
				Aliases: []string{"r"},
				Usage:   "full directory tree will be returned",
			},
			&cli.BoolFlag{
				Name:  "showdeleted",
				Usage: "deleted files and folders that can be undeleted will be displayed",
			},
			&cli.BoolFlag{
				Name:  "nofiles",
				Usage: "only the folder (sub)structure will be returned",
			},
			&cli.BoolFlag{
				Name:  "noshares",
				Usage: "only user's own folders and files will be displayed",
			},
		},
	},
	{
		Name:    "up",
		Aliases: []string{"uploadfile"},
		Usage:   "Upload a file.",
		Action: func(c *cli.Context) error {
			auth := c.Context.Value("auth").(string)
			src := c.Args().First()
			file, err := os.Open(src)
			if err != nil {
				return printError(fmt.Errorf("cannot open %s", src))
			}
			defer file.Close()

			dest := c.Args().Get(1)
			if err := uploadfile(auth, dest, file); err != nil {
				return printError(err)
			}
			return nil
		},
		ArgsUsage: "<SOURCE_FILE> <DEST_PATH>",
	},
}

func main() {
	os.Exit(run())
}

func run() int {
	app := cli.NewApp()
	app.Name = APP_NAME
	app.Version = VERSION
	app.Before = login
	app.Commands = commands
	if err := app.Run(os.Args); err != nil {
		return 1
	} else {
		return 0
	}
}

func login(c *cli.Context) error {
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

	//lint:ignore SA1029 set auth
	c.Context = context.WithValue(c.Context, "auth", *auth)
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

func btoi(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func printError(err error) error {
	red := color.New(color.FgRed).FprintfFunc()
	red(os.Stderr, "%v\n", err)
	return err
}
