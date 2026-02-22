package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/t-chov/pcloud-uploader/pcloud"
	"github.com/urfave/cli/v2"
)

var printFolder = color.New(color.FgBlue).PrintlnFunc()

func run() int {
	client := pcloud.NewClient("")
	return runWithClient(client)
}

func runWithClient(client pcloud.API) int {
	app := cli.NewApp()
	app.Name = APP_NAME
	app.Version = VERSION
	app.Before = func(c *cli.Context) error {
		return login(c, client)
	}
	app.Commands = buildCommands(client)
	if err := app.Run(os.Args); err != nil {
		return 1
	}
	return 0
}

func buildCommands(client pcloud.API) []*cli.Command {
	return []*cli.Command{
		{
			Name:    "ls",
			Aliases: []string{"listfolder"},
			Usage:   "Receive data for a folder.",
			Action: func(c *cli.Context) error {
				auth := c.Context.Value("auth").(string)
				opts := pcloud.ListFolderOptions{
					Recursive:   c.Bool("recursive"),
					ShowDeleted: c.Bool("showdeleted"),
					NoFiles:     c.Bool("nofiles"),
					NoShares:    c.Bool("noshares"),
				}
				result, err := client.ListFolder(auth, c.Args().First(), opts)
				if err != nil {
					return printError(err)
				}
				printListfolder("/", result.Metadata)
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
				result, err := client.UploadFile(auth, dest, file)
				if err != nil {
					return printError(err)
				}
				for _, metadata := range result.Metadata {
					fmt.Println(metadata.Hash.String())
				}
				return nil
			},
			ArgsUsage: "<SOURCE_FILE> <DEST_PATH>",
		},
	}
}

func login(c *cli.Context, client pcloud.API) error {
	username, err := loadFromInput(ENV_PCLOUD_USERNAME, "username")
	if err != nil {
		return err
	}
	password, err := loadFromInput(ENV_PCLOUD_PASSWORD, "password")
	if err != nil {
		return err
	}

	auth, err := client.Authenticate(username, password)
	if err != nil {
		return err
	}
	if auth == "" {
		return fmt.Errorf("empty auth")
	}

	//lint:ignore SA1029 set auth
	c.Context = context.WithValue(c.Context, "auth", auth)
	return nil
}

func loadFromInput(envName string, title string) (string, error) {
	value := os.Getenv(envName)
	if value == "" {
		fmt.Printf("%s: ", title)
		if _, err := fmt.Scan(&value); err != nil {
			return "", fmt.Errorf("scan %s: %v", title, err)
		}
	}
	value = strings.TrimSpace(value)
	return value, nil
}

func printListfolder(path string, ls pcloud.FolderMetadata) {
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

func printError(err error) error {
	red := color.New(color.FgRed).FprintfFunc()
	red(os.Stderr, "%v\n", err)
	return err
}
